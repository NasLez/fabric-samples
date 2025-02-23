package test

import (
	"context"
	"fabric_ebl/domain"
	"fabric_ebl/service"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ebl"
	"testing"
)

type mockConnectService struct {
	mock.Mock
}

func (m *mockConnectService) Contract(ctx context.Context, userName string, companyName string) (contract *gateway.Contract, closeFunc service.CloseFunc, err error) {
	args := m.Called(ctx, userName, companyName)
	return args.Get(0).(*gateway.Contract), args.Get(1).(service.CloseFunc), args.Error(2)
}

func (m *mockConnectService) ParseToken(ctx context.Context, req string) (token string, claims map[string]any, userId int64, user *domain.UserDO, company *domain.CompanyDO, err error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Get(1).(map[string]any), args.Get(2).(int64), args.Get(3).(*domain.UserDO), args.Get(4).(*domain.CompanyDO), args.Error(5)
}

func TestFabricEblServiceImplQueryAllEblList(t *testing.T) {
	// 创建模拟的 ConnectService
	mockConnectService := new(mockConnectService)

	// 创建 FabricEblServiceImpl 实例，并注入模拟的 ConnectService
	service := service.FabricEblServiceImpl{
		P: service.Param{
			FabricEblRepo:  nil,
			ConnectService: mockConnectService,
		},
	}

	// 模拟 ParseToken 的返回值
	token := "test_token"
	claims := map[string]any{
		"user_id": 123,
	}
	userId := int64(123)
	user := &domain.UserDO{
		ID:   123,
		Name: "test_user",
	}
	company := &domain.CompanyDO{
		ID:   456,
		Name: "test_company",
	}
	mockConnectService.On("ParseToken", context.Background(), token).Return(token, claims, userId, user, company, nil)

	// 模拟 Contract 的返回值
	contract := &gateway.Contract{}
	closeFunc := func() {}
	mockConnectService.On("Contract", context.Background(), user.Name, company.Name).Return(contract, closeFunc, nil)

	// 模拟 SubmitTransaction 的返回值
	result := []byte("test_result")
	mockConnectService.On("SubmitTransaction", contract, "GetEblByRangeWithPagination", "", "", "").Return(result, nil)

	// 创建测试请求
	req := &fabric_ebl.QueryAllEblListReq{
		Token: token,
	}

	// 调用目标函数
	resp, err := service.QueryAllEblList(context.Background(), req)

	// 断言结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(resp.EblList), 0)
	assert.Equal(t, resp.Bookmark, "")
	assert.Equal(t, resp.FetchedRecordsCount, int64(0))

	// 验证模拟的方法调用
	mockConnectService.AssertExpectations(t)
}
