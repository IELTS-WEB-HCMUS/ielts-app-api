- Với phần source BE (ngôn ngữ Golang - framework GIN): Nhóm đã xây dựng source dựa trên **Clean Architecture** với các thư mục được tổ chức rõ ràng, tách biệt trách nhiệm và phụ thuộc, giúp dễ bảo trì và mở rộng. Hơn nữa, với phần source BE nhóm tụi em đã deploy API lên server của Render. Tuy nhiên, Render sẽ cho hệ thống API vào trạng thái sleep nếu như sau 15' mà không có bất kì 1 API nào được request đến. Và nếu như server API bị sleep thì lần gọi API tiếp theo sẽ phải mất 1 khoảng thời gian khá lâu để gọi được API (>=50s) nên nhóm em đã có sử dụng thêm 1 service cron job free có tên là [cronjob.org] (https://cron-job.org/en/) để cứ mỗi 10', nhóm em sẽ cho con cron job này gọi đến 1 API Get bất kì để giữ server không bị sleep. 
  1. **cmd**: Thư mục này thường chứa các điểm khởi đầu chính của ứng dụng. Trong trường hợp này, có thể có một tệp chính trong `cmd` để khởi chạy ứng dụng.
  2. **common**: Thư mục này chứa các đoạn mã tiện ích và các chức năng dùng chung trong toàn bộ ứng dụng:
       - `db.go`: Chứa cấu hình và các hàm kết nối cơ sở dữ liệu.
       - `error.go`, `error_messages.go`: Xử lý các loại lỗi tùy chỉnh và thông điệp lỗi.
       - `helper.go`, `utils.go`: Các hàm trợ giúp chung có thể tái sử dụng trong các package khác.
       - `jwt.go`: Quản lý mã JWT, xử lý mã hóa và giải mã để xác thực.
  3. **config**: Thư mục này có thể chứa các thiết lập và cấu hình liên quan đến ứng dụng, như các tệp cấu hình hoặc các tham số cấu hình hệ thống.
  4. **internal**: Thư mục `internal` chứa các thành phần cốt lõi của ứng dụng, không được truy cập từ bên ngoài, và được chia thành các thư mục con theo vai trò:
     - handlers: chứa các tệp xử lý logic API như `authen.go` (Xử lý logic xác thực người dùng), `check_health.go` (Kiểm tra tình trạng của dịch vụ), `target.go`, `user.go`(Xử lý các chức năng liên quan đến đối tượng `target` và `user`)
     - models: chứa các tệp định nghĩa cấu trúc dữ liệu của ứng dụng (models) như `authen.go`, `user.go`, `target.go`, `role.go`, `base.go` để định nghĩa cấu trúc dữ liệu cho các đối tượng liên quan như `user`, `role` (vai trò), `target` (đối tượng), v.v.
     - repositories: chứa các tệp làm việc với cơ sở dữ liệu như `authen.go`, `user.go`, `target.go`, `base.go` để xử lý các thao tác CRUD (tạo, đọc, cập nhật, xóa) cho từng đối tượng cụ thể, đảm bảo tách biệt việc truy cập dữ liệu khỏi logic nghiệp vụ.
     - services: chứa các tệp xử lý logic nghiệp vụ như `authen.go`, `user.go`, `target.go`, `base.go` để xử lý logic nghiệp vụ cho từng đối tượng cụ thể, đóng vai trò trung gian giữa các `handler` và `repository`.
     - middleware: chứa các tệp xử lý trung gian (middleware) như  `admin_authentication.go`, `user_authentication.go` để quản lý xác thực và phân quyền, kiểm tra quyền truy cập của admin và user.
  5. **pkg**: để chứa các package bổ trợ như `postgres` (Chứa các tiện ích hoặc cấu hình dành riêng cho PostgreSQL) hay `utils` (Các tiện ích khác có thể cần thiết cho nhiều phần của ứng dụng).
  6. **script**: Chứa các script hoặc tệp SQL cần thiết cho ứng dụng
  7. **Các tệp cấu hình chung**:
  - `.env`: Tệp môi trường, chứa các biến môi trường như cấu hình cơ sở dữ liệu, thông tin API, v.v.
  - `go.mod`, `go.sum`: Quản lý các module và dependencies của Go.
  - `.gitignore`: Xác định các tệp và thư mục không cần thiết đẩy lên Git.
 
Cấu trúc này giúp phân chia rõ ràng giữa các tầng (layer) của ứng dụng:
- **handlers** cho việc xử lý HTTP.
- **services** cho logic nghiệp vụ.
- **repositories** cho truy cập dữ liệu.
- **models** cho các định nghĩa dữ liệu.
