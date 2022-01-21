# 精确确认BAC边界
## 1 提取所有与载体末端重叠的reads
## 2 截断reads仅保留HindIII位点下游的序列
## 2 提取基因组HindIII酶切位点下游150 bp序列
## 4 将截断reads在基因组序列中查找，确定具体的坐标

### 2021-12-31 更新
1. 修改plateID参数解析方式，将原本逗号分隔的数组指定改为指定单元编号，及SP01、SP02、SP03等
2. 修改runBowtie2.go line 20, "100" to string(threads)
3. 运行SP03
    ```bash
   ~/software/toolkit/BACtoPanda/BACtoPanda 
        -genome=/home/Xuwei/reference_genome/asm.cleaned.fasta.assembly.assembly.FINAL.fasta 
        -sourceDir=/mnt/raw_data/data_shenfujun/DNA/X101SC21083152-Z01-J004/2.cleandata/
        -unit=SP03 
        -threads=40  
        -profile=/home/Xuwei/reference_genome/EnzymeProfile/panda-genome.bed
   