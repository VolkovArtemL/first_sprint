package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("ожидается 3 части, получено %d", len(parts))
	}

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка конвертации шагов '%s': %w", parts[0], err)
	}

	activity := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга длительности '%s': %w", parts[2], err)
	}

	return steps, activity, duration, nil

}

func distance(steps int, height float64) float64 {

	stepsLength := stepLengthCoefficient * height

	distanceM := stepsLength * float64(steps)

	return distanceM / mInKm

}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)

	durationHours := duration.Hours()

	return dist / durationHours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println("Ошибка при парсинге данных:", err)
		return "", err
	}

	var distanceKm, speed, calories float64
	var activityType string

	// Определяем тип тренировки
	switch activity {
	case "Ходьба":
		activityType = "Ходьба"
		distanceKm = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)

	case "Бег":
		activityType = "Бег"
		distanceKm = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)

	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}

	// Формируем строку с результатом
	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activityType,
		duration.Hours(),
		distanceKm,
		speed,
		calories)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	switch {
	case steps <= 0:
		return 0, fmt.Errorf("некорректное количество шагов: %d (должно быть > 0)", steps)
	case weight <= 0:
		return 0, fmt.Errorf("некорректный вес: %.2f кг (должен быть > 0)", weight)
	case height <= 0:
		return 0, fmt.Errorf("некорректный рост: %.2f м (должен быть > 0)", height)
	case height > 3:
		return 0, fmt.Errorf("нереалистичный рост: %.2f м (максимальный рост 3 м)", height)
	case weight > 300:
		return 0, fmt.Errorf("нереалистичный вес: %.2f кг (максимальный вес 300 кг)", weight)
	case duration <= 0:
		return 0, fmt.Errorf("некорректная продолжительность: %v (должна быть > 0)", duration)
	case duration > 24*time.Hour:
		return 0, fmt.Errorf("продолжительность не может превышать 24 часа: %v", duration)
	}

	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, fmt.Errorf("ошибка расчёта скорости: скорость = %.2f км/ч", speed)
	}
	if speed > 30 {

		return 0, fmt.Errorf("нереалистичная скорость: %.2f км/ч", speed)
	}

	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH
	if calories < 0 {
		return 0, fmt.Errorf("ошибка расчёта: получено отрицательное значение калорий %.2f", calories)
	}

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	switch {
	case steps <= 0:
		return 0, fmt.Errorf("некорректное количество шагов: %d (должно быть > 0)", steps)
	case weight <= 0:
		return 0, fmt.Errorf("некорректный вес: %.2f кг (должен быть > 0)", weight)
	case height <= 0:
		return 0, fmt.Errorf("некорректный рост: %.2f м (должен быть > 0)", height)
	case height > 3:
		return 0, fmt.Errorf("нереалистичный рост: %.2f м (максимальный рост 3 м)", height)
	case weight > 300:
		return 0, fmt.Errorf("нереалистичный вес: %.2f кг (максимальный вес 300 кг)", weight)
	case duration <= 0:
		return 0, fmt.Errorf("некорректная продолжительность: %v (должна быть > 0)", duration)
	case duration > 24*time.Hour:
		return 0, fmt.Errorf("продолжительность не может превышать 24 часа: %v", duration)
	}

	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, fmt.Errorf("ошибка расчёта скорости: скорость = %.2f км/ч", speed)
	}
	if speed > 15 {

		return 0, fmt.Errorf("нереалистичная скорость: %.2f км/ч", speed)
	}

	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH
	if calories < 0 {
		return 0, fmt.Errorf("ошибка расчёта: получено отрицательное значение калорий %.2f", calories)
	}

	walkingCalories := calories * walkingCaloriesCoefficient

	return walkingCalories, nil

}
