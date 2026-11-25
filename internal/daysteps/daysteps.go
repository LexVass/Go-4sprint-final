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
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("Неправильное количество данных: %d из необходимых 2", len(parts))
	}
	stepsCount, err1 := strconv.Atoi(parts[0])
	if err1 != nil {
		return 0, 0, err1
	}
	if stepsCount <= 0 {
		return 0, 0, fmt.Errorf("Количество шагов %d меньше или равно 0", stepsCount)
	}
	dur, err2 := time.ParseDuration(parts[1])
	if err2 != nil {
		return 0, 0, err2
	}
	return stepsCount, dur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	stepsCount, durTime, err1 := parsePackage(data)
	if err1 != nil {
		fmt.Println(err1)
		return ""
	}
	distantionKm := float64(stepsCount) * stepLength / mInKm
	kcal, err2 := spentcalories.WalkingSpentCalories(stepsCount, weight, height, durTime)
	if err2 != nil {
		fmt.Println(err2)
		return ""
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", stepsCount, distantionKm, kcal)
}
