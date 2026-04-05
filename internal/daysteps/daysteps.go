package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	parts := strings.Split(data, ", ")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("ожидается 2 значения с разделением, получено %d", len(parts))
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка конвертации: %w", err)
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("Количество шагов должно быть больше нуля")
	}

	durationWalk, err := time.ParseDuration(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга длительности '%s': %w", parts[1], err)
	}

	return steps, durationWalk, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return ""
	}

	if steps <= 0 {
		fmt.Println("Количество шагов должно быть больше 0")
		return ""
	}

	distanceMeters := float64(steps) * stepLength

	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %2.f км.\nВы сожгли %2.f ккал. ", steps, distanceKm, calories)

	return result
}
