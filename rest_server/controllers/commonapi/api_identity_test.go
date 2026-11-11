package commonapi

import (
	"testing"

	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/LumiWave/inno-point-manager/rest_server/model"
)

func testAppPoints() map[int64]*model.AppPointInfo {
	return map[int64]*model.AppPointInfo{
		1: {PointsMap: map[int64]*model.PointInfo{10: {}}},
		2: {PointsMap: map[int64]*model.PointInfo{20: {}, 21: {}}},
		3: nil,
	}
}

func TestAppHasPointIn(t *testing.T) {
	appPoints := testAppPoints()
	cases := []struct {
		appID, pointID int64
		want           bool
	}{
		{1, 10, true},
		{2, 21, true},
		{1, 20, false}, // 다른 앱의 포인트
		{2, 10, false},
		{3, 10, false}, // 메타가 비어 있는 앱
		{9, 10, false}, // 없는 앱
		{0, 10, false}, // 앱 정보 없는 토큰
	}
	for _, c := range cases {
		if got := appHasPointIn(appPoints, c.appID, c.pointID); got != c.want {
			t.Errorf("app %d point %d : %v, 기대 %v", c.appID, c.pointID, got, c.want)
		}
	}
}

func TestPointsBelongToApp(t *testing.T) {
	appPoints := testAppPoints()
	member := []*context.Point{{PointID: 20}, nil, {PointID: 21}}

	if !pointsBelongToApp(appPoints, member, 2) {
		t.Error("자기 앱 회원의 포인트를 다른 앱 것으로 판단함")
	}
	if pointsBelongToApp(appPoints, member, 1) {
		t.Error("다른 앱 회원의 포인트 조회가 통과함")
	}
	if pointsBelongToApp(appPoints, nil, 1) {
		t.Error("포인트가 없는데 소속으로 판단함")
	}
}
