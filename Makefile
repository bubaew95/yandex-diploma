mock:
	mockgen -source=./internal/core/ports/order.go -destination=./internal/infra/repository/mock/order.go -package=mock -mock_names=OrderRepository=MockOrderRepository \
	&& mockgen -source=./internal/core/ports/user.go -destination=./internal/infra/repository/mock/user.go -package=mock -mock_names=UserRepository=MockUserRepository