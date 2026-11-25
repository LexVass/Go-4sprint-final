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
	//lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("Неправильное количество данных: %d из необходимых 3", len(parts))
	}
	stepsCount, err1 := strconv.Atoi(parts[0])
	if err1 != nil {
		return 0, "", 0, err1
	}
	if stepsCount <= 0 {
		return 0, "", 0, fmt.Errorf("Количество шагов (%d) должно быть больше нуля", stepsCount)
	}
	dur, err2 := time.ParseDuration(parts[2])
	if err2 != nil {
		return 0, "", 0, err2
	}
	if dur <= 0 {
		return 0, "", 0, fmt.Errorf("Продолжительность активности (%v) должна быть больше нуля", dur)
	}
	return stepsCount, parts[1], dur, nil
}

func distance(steps int, height float64) float64 {
	stepLenght := height * stepLengthCoefficient
	return stepLenght * float64(steps) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distKm := distance(steps, height)
	return distKm / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	stepsCount, typeTraining, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	switch typeTraining {

	case "Ходьба":
		caloriesTraining, err := WalkingSpentCalories(stepsCount, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"Тип треннировки: %s\n"+
				"Длительность: %.2f\n ч."+
				"Дистанция: %.2f\n км."+
				"Скорость: %.2f\n км/ч"+
				"Сожгли калорий: %.2f\n",
			typeTraining,
			duration.Hours(),
			distance(stepsCount, height),
			meanSpeed(stepsCount, height, duration),
			caloriesTraining), nil

	case "Бег":
		caloriesTraining, err := RunningSpentCalories(stepsCount, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"Тип треннировки: %s\n"+
				"Длительность: %.2f ч.\n"+
				"Дистанция: %.2f км.\n"+
				"Скорость: %.2f км/ч\n"+
				"Сожгли калорий: %.2f\n",
			typeTraining,
			duration.Hours(),
			distance(stepsCount, height),
			meanSpeed(stepsCount, height, duration),
			caloriesTraining), nil

	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("Количество шагов (%d) должно быть больше нуля", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("Масса (%.2f) должна быть больше нуля", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("Высота (%.2f) должна быть больше нуля", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("Продолжительность тренировки (%v) должна быть больше нуля", duration)
	}
	meanSpeed := meanSpeed(steps, height, duration)
	return (weight * meanSpeed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("Количество шагов (%d) должно быть больше нуля", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("Масса (%.2f) должна быть больше нуля", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("Высота (%.2f) должна быть больше нуля", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("Продолжительность ходьбы (%v) должна быть больше нуля", duration)
	}
	meanSpeed := meanSpeed(steps, height, duration)
	return ((weight * meanSpeed * duration.Minutes()) / minInH) * walkingCaloriesCoefficient, nil
}
