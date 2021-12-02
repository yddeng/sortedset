package grank

import (
	"github.com/yddeng/grank/skiplist"
	"time"
)

type LeaderboardOrder int

const (
	OrderDeclining LeaderboardOrder = 0 // 降序
	OrderAscending LeaderboardOrder = 1 // 升序
)

type LeaderboardUpdateStrategy int

const (
	UpdateStrategyLast   LeaderboardUpdateStrategy = 0 // 最新成绩
	UpdateStrategyBetter LeaderboardUpdateStrategy = 1 // 最好成绩
	UpdateStrategySum    LeaderboardUpdateStrategy = 2 // 总成绩
)

type Leaderboard struct {
	StatisticName  string                    // 排行榜名
	Order          LeaderboardOrder          // 排序
	UpdateStrategy LeaderboardUpdateStrategy // 更新策略
	Version        int                       // 版本号
	CreateAt       int64                     // 创建时间
	NextResetAt    int64                     // 下一次重置时间
	Period         time.Duration             // 版本持续周期
	skipList       *skiplist.SkipList
}

type Statistic struct {
	User  string
	Score float64
}

func New(statisticName string, order LeaderboardOrder, updateStrategy LeaderboardUpdateStrategy) *Leaderboard {
	return &Leaderboard{
		StatisticName:  statisticName,
		Order:          order,
		UpdateStrategy: updateStrategy,
		Version:        1,
		CreateAt:       time.Now().Unix(),
		NextResetAt:    0,
		Period:         0,
	}
}

func (this *Leaderboard) UpdateStatistic(key string, score float64) bool {
	this.skipList.Remove()
}
