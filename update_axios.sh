#!/bin/bash

# Update import statements
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/import axios from '\''@\/utils\/axios'\'';/import axiosInstance from '\''@\/utils\/axios'\'';/g' {} \;

# Update axios.get calls with proper /api prefix
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.get<[^>]*>('"'"'api\//await axiosInstance.get<\1>('\'\/api\//g' {} \;
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.get<[^>]*>(`api\//await axiosInstance.get<\1>(`\/api\//g' {} \;
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.get<[^>]*>('"'"'api/await axiosInstance.get<\1>('\'\/api/g' {} \;
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.get<[^>]*>(`api/await axiosInstance.get<\1>(`\/api/g' {} \;

# Update axios.post calls with proper /api prefix
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.post('"'"'api\//await axiosInstance.post('\'\/api\//g' {} \;
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.post(`api\//await axiosInstance.post(`\/api\//g' {} \;
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.post('"'"'api/await axiosInstance.post('\'\/api/g' {} \;
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.post(`api/await axiosInstance.post(`\/api/g' {} \;

# Update axios.patch calls with proper /api prefix
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.patch('"'"'api\//await axiosInstance.patch('\'\/api\//g' {} \;
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.patch(`api\//await axiosInstance.patch(`\/api\//g' {} \;
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.patch('"'"'api/await axiosInstance.patch('\'\/api/g' {} \;
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.patch(`api/await axiosInstance.patch(`\/api/g' {} \;

# Update axios.delete calls with proper /api prefix
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.delete('"'"'api\//await axiosInstance.delete('\'\/api\//g' {} \;
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.delete(`api\//await axiosInstance.delete(`\/api\//g' {} \;
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.delete('"'"'api/await axiosInstance.delete('\'\/api/g' {} \;
find /home/ihor/Desktop/projects/cows-shelter/frontend/src/store/slices -type f -name "*.ts" -exec sed -i 's/await axios\.delete(`api/await axiosInstance.delete(`\/api/g' {} \;
