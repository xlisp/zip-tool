#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
将 part000.txt 到 part757.txt 分配到 7 个目录中
每个目录大约包含 108-109 个文件
"""

import os
import shutil
import sys

def create_directories_and_distribute_files():
    """创建7个目录并分配文件"""
    
    # 总文件数和目录数
    total_files = 758  # part000.txt 到 part757.txt
    num_directories = 7
    
    # 计算每个目录应该包含的文件数
    files_per_dir = total_files // num_directories
    remaining_files = total_files % num_directories
    
    print(f"总文件数: {total_files}")
    print(f"目录数: {num_directories}")
    print(f"每个目录基本文件数: {files_per_dir}")
    print(f"剩余文件数: {remaining_files}")
    print()
    
    # 创建目录名列表
    directories = [f"dir_{i+1:02d}" for i in range(num_directories)]
    
    # 创建目录
    for directory in directories:
        if not os.path.exists(directory):
            os.makedirs(directory)
            print(f"创建目录: {directory}")
    
    print()
    
    # 分配文件
    current_file_index = 0
    
    for dir_index, directory in enumerate(directories):
        # 前几个目录如果有剩余文件，多分配一个
        current_dir_file_count = files_per_dir + (1 if dir_index < remaining_files else 0)
        
        print(f"目录 {directory} 将包含 {current_dir_file_count} 个文件:")
        
        for i in range(current_dir_file_count):
            file_number = current_file_index
            source_file = f"part{file_number:03d}.txt"
            destination = os.path.join(directory, source_file)
            
            # 检查源文件是否存在
            if os.path.exists(source_file):
                # 移动文件（如果要复制而不是移动，使用 shutil.copy2）
                try:
                    shutil.move(source_file, destination)
                    print(f"  移动: {source_file} -> {destination}")
                except Exception as e:
                    print(f"  错误: 无法移动 {source_file}: {e}")
            else:
                print(f"  警告: 文件 {source_file} 不存在")
            
            current_file_index += 1
        
        print()
    
    print("文件分配完成!")

def show_distribution_plan():
    """显示分配计划（不实际移动文件）"""
    total_files = 758
    num_directories = 7
    files_per_dir = total_files // num_directories
    remaining_files = total_files % num_directories
    
    print("=== 文件分配计划 ===")
    print(f"总文件数: {total_files}")
    print(f"目录数: {num_directories}")
    print()
    
    current_file_index = 0
    
    for dir_index in range(num_directories):
        directory = f"dir_{dir_index+1:02d}"
        current_dir_file_count = files_per_dir + (1 if dir_index < remaining_files else 0)
        
        start_file = current_file_index
        end_file = current_file_index + current_dir_file_count - 1
        
        print(f"{directory}: part{start_file:03d}.txt ~ part{end_file:03d}.txt ({current_dir_file_count} 个文件)")
        
        current_file_index += current_dir_file_count
    
    print()

if __name__ == "__main__":
    print("文件分配脚本")
    print("=" * 50)
    
    # 显示分配计划
    show_distribution_plan()
    
    # 询问用户是否继续
    if len(sys.argv) > 1 and sys.argv[1] == "--execute":
        create_directories_and_distribute_files()
    else:
        response = input("是否要执行文件移动？(y/N): ").strip().lower()
        if response in ['y', 'yes', '是']:
            create_directories_and_distribute_files()
        else:
            print("取消操作。如要直接执行，可使用参数 --execute")

