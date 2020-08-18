package main

/*
	Вариант решения задачи нахождения размера максимального блока, 
	при котором сохраняется только размер максимального блока.

	Максимальный размер массива для Win32 при 2 цветах равен около 75 млн ячеек.

	При увеличении количества цветов размер массива почти пропорционально падает
	по непонятной причине!
*/

import (
	"fmt"
	"math/rand"
	"time"
	"testing"
)

type cell struct {
	color int	// цвет ячейки
	block_id int	// № блока, в который вошла ячейка
}

type block struct {
	id int 		// номер блока по порядку создания
	size int	// размер блока
}

// основные константы
const (
        row_count = 10000	// количество строк
	col_count = 10000        // количество столбцов
	color_count = 2         // число разных цветов 
)

// глобальные переменные программы
var (
	max_block_size int = 1 		// размер максимального блока (по умолчанию = 1)
	
	cells [row_count][col_count]*cell 	// матрица ячеек
	
	processed_cells int		// число обработанных ячеек (для контроля)

	blocks_count int		// текущее число блоков

	blocks []*block 				// список блоков
)

/*
	Функция добавления ячейки в состав блока
*/
func (b *block) add_cell(c *cell) {

        processed_cells = processed_cells + 1

	check_max := func () { // вложенная функция проверки максимума

		// если при соединении блоков образовался самый большой блок на поле, то это надо зафиксировать 
		if b.size > max_block_size {
			max_block_size = b.size
		}  
	}

   	if (c.block_id == b.id) {   // если ячейка уже входит в этот блок, то сразу на выход
		return
  	}

   	if c.block_id >= 0  {
		if c.block_id < b.id {
			return // если ячейка уже относится к более старому блоку, то ничего не делаем
		} else { 
  			// если ячейка относится к более молодому блоку, то надо весь тот блок присоединить к текущему
 			b2 := blocks[c.block_id]
			b.size = b.size + b2.size
			b2.size = 0
		}
  	} else {

	 	// добавление независимой ячейки в текущий блок
		b.size = b.size + 1 	  		// размер текущего блока увеличивается на 1

		c.block_id = b.id			// ячейка приписывается к блоку
	}	

 	check_max()
	return

}


/*
  Просто функция определения минимума из двух чисел
*/
func min(x, y int) int {
 if x < y {
   return x
 }
 return y
}

/*
	Основная функция приложения
*/
func main() {

	t0 := time.Now()

//	minimal_dimension := min(row_count, col_count)

	var i int = 0
	var j int = 0

	r := rand.New(rand.NewSource(10))

        // инициализируем ячейки в стартовое состояние
	for i = 0; i < row_count; i++ {
		for j = 0; j < col_count; j++ {
			c := new(cell)
			c.color = r.Intn(color_count)
			c.block_id = -2 // по умолчанию блок не задан
			cells[i][j] = c
		}
	}

	fmt.Println("Matrix prepared")

        // основной цикл подсчета блоков
	for i = 0; i < row_count; i++ {
		for j = 0; j < col_count; j++ {

			// соединение с предыдущим рядом по вертикали, начиная со второго ряда
			if i > 0 && cells[i][j].color == cells[i-1][j].color && cells[i][j].block_id < 0 {
				b := blocks[cells[i-1][j].block_id]
				b.add_cell(cells[i][j])
				continue
			}

                        // соединение с предыдущей ячейкой по горизонтали в случае если текущая ячейка еще не присоединена к блоку
			if j > 0 && cells[i][j].color == cells[i][j-1].color && cells[i][j].block_id < 0 {
				b := blocks[cells[i][j-1].block_id]
				b.add_cell(cells[i][j])
				continue
			}

                        // соединение с предыдущей ячейкой по горизонтали в случае если текущая ячейка УЖЕ присоединена к блоку
			if j > 0 && cells[i][j].color == cells[i][j-1].color && cells[i][j].block_id >= 0 && cells[i][j].block_id != cells[i][j-1].block_id {
  				if cells[i][j].block_id < cells[i][j-1].block_id {
					b := blocks[cells[i][j].block_id]
					b.add_cell(cells[i][j-1])
				} else {
					b := blocks[cells[i][j-1].block_id]
					b.add_cell(cells[i][j])
				}
				continue
			}

 			// если ячейка не присоединилась к существующему блоку ни слева ни сверху, то создаём новый блок и добавляем в него эту ячейку
			// эти же действия происходит в самой первой ячейке
			if cells[i][j].block_id < 0 {
				var b block
				b.id = blocks_count   	// номер по порядку создания
				blocks_count ++
				blocks = append(blocks, &b)
				b.add_cell(cells[i][j])
			}
			
			// на каждой итерации проверяем, что обще число активных блоков не превысило размер меньшего из измерений матрицы,
			// так как никакой самый хитрый блок не может продолжаться если после него создано еше minimal_dimension новых блоков
			// TO-DO - написанный нижде код неработтоспособен - так как просто вырезает все блоки кроме последних,
			// а надо ещё дополнительно сортировать их по мере добавления ячеек - т.е. должны остаться не просто N посоедних блоков,
			// а N блоков, в которые последними добавлялись ячейки 	
			/*
			if len(blocks) > minimal_dimension {
				blocks = blocks[minimal_dimension:] // оставляем только последних N блоков в живых
			}

			*/
		}
	}

	fmt.Println("Processed cells count is", processed_cells)
	fmt.Println("Total blocks count is", len(blocks))
	fmt.Println("Max block size is", max_block_size)

	blocks = blocks[:0]

	fmt.Println("----------------------") 

	t1 := time.Now()
	fmt.Println("Benchmark", t1.Sub(t0))

        // ожидание Enter для выхода из программы
	fmt.Print("Press Enter to exit")
	var s string
	fmt.Scanln(&s)


}

func BenchmarkMain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main()
	}
}
