package controllers

import (
	"context"
	"example/config"
	"example/entities"
	"example/tracing"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var tracer = tracing.Tracer

func recordMetrics(ctx context.Context, path string, start time.Time, status int) {
	duration := time.Since(start).Seconds()

	tracing.RequestCount.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("http.path", path),
			attribute.Int("http.status", status),
		),
	)

	tracing.RequestDuration.Record(ctx, duration,
		metric.WithAttributes(
			attribute.String("http.path", path),
			attribute.Int("http.status", status),
		),
	)
}

func GetCategories(c *fiber.Ctx) error {
	ctx, span := tracer.Start(c.UserContext(), "GetCategories")
	defer span.End()

	start := time.Now()
	var categories []entities.Category

	if err := config.GormDB.Find(&categories).Error; err != nil {
		span.RecordError(err)
		recordMetrics(ctx, c.Path(), start, http.StatusInternalServerError)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	recordMetrics(ctx, c.Path(), start, http.StatusOK)
	return c.JSON(categories)
}

func GetCategory(c *fiber.Ctx) error {
	ctx, span := tracer.Start(c.UserContext(), "GetCategory")
	defer span.End()

	start := time.Now()

	id := c.Params("id")
	span.SetAttributes(attribute.String("category.id", id))
	var category entities.Category

	if err := config.GormDB.First(&category, id).Error; err != nil {
		span.RecordError(err)
		recordMetrics(ctx, c.Path(), start, http.StatusNotFound)
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	recordMetrics(ctx, c.Path(), start, http.StatusOK)
	return c.JSON(category)
}

func CreateCategory(c *fiber.Ctx) error {
	ctx, span := tracer.Start(c.UserContext(), "CreateCategory")
	defer span.End()

	start := time.Now()
	var category entities.Category

	if err := c.BodyParser(&category); err != nil {
		span.RecordError(err)
		recordMetrics(ctx, c.Path(), start, http.StatusBadRequest)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := config.GormDB.Create(&category).Error; err != nil {
		span.RecordError(err)
		recordMetrics(ctx, c.Path(), start, http.StatusInternalServerError)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	recordMetrics(ctx, c.Path(), start, http.StatusCreated)
	return c.Status(201).JSON(category)
}

func UpdateCategory(c *fiber.Ctx) error {
	ctx, span := tracer.Start(c.UserContext(), "UpdateCategory")
	defer span.End()

	start := time.Now()
	id := c.Params("id")
	span.SetAttributes(attribute.String("category.id", id))

	var category entities.Category
	if err := config.GormDB.First(&category, id).Error; err != nil {
		span.RecordError(err)
		recordMetrics(ctx, c.Path(), start, http.StatusNotFound)
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	var input entities.Category
	if err := c.BodyParser(&input); err != nil {
		span.RecordError(err)
		recordMetrics(ctx, c.Path(), start, http.StatusBadRequest)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	category.Name = input.Name
	if err := config.GormDB.Save(&category).Error; err != nil {
		span.RecordError(err)
		recordMetrics(ctx, c.Path(), start, http.StatusInternalServerError)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	recordMetrics(ctx, c.Path(), start, http.StatusOK)
	return c.JSON(category)
}

func DeleteCategory(c *fiber.Ctx) error {
	ctx, span := tracer.Start(c.UserContext(), "DeleteCategory")
	defer span.End()

	start := time.Now()
	id := c.Params("id")
	span.SetAttributes(attribute.String("category.id", id))

	var category entities.Category
	if err := config.GormDB.First(&category, id).Error; err != nil {
		span.RecordError(err)
		recordMetrics(ctx, c.Path(), start, http.StatusNotFound)
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	if err := config.GormDB.Delete(&category).Error; err != nil {
		span.RecordError(err)
		recordMetrics(ctx, c.Path(), start, http.StatusInternalServerError)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	recordMetrics(ctx, c.Path(), start, http.StatusNoContent)
	return c.SendStatus(204)
}
