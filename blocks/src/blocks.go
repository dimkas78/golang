package main

import (
	"fmt"
	"math/rand"
)

type cell struct {
        id int		// № ячейки - порядковый номер, начиная с левого верхнего угла
	color int	// цвет ячейки
	row int         // номер строки (считаем сверху)
	col int         // номер столбца (считаем слева)
	block_id int	// № блока, в который вошла ячейка
}

type block struct {
	id int 		// номер блока по порядку создания
	size int	// размер блока
	start *cell	// ячейка, в которой этот блок стартовал
	cells []*cell   // список указателей на ячейки, по порядку вхождения в блок
	parent_id int	// номер блока, который "поглотил" этот блок при слиянии блоков
}

// основные константы
const (
        row_count = 40 	// количество строк
	col_count = 90        // количество столбцов
	color_count = 2         // число разных цветов 
)

// глобальные переменные программы
var (
	max_block_size int = 1 		// размер максимального блока (по умолчанию = 1)
	max_block_id int = -1 		// № максимального блока (по умолчанию - не определен)
	
	cells [row_count][col_count]*cell 	// матрица ячеек

	blocks []*block 				// список блоков
)


/*
	Функция инициализации ячейки 
*/
func (c *cell) init_cell(in_row int, in_col int, in_color int) {

	c.row = in_row
	c.col = in_col
	c.id = in_row * col_count + in_col	// номер по порядку от левого верхнего угла

	c.color = in_color

	c.block_id = max_block_id - 1 // по умолчанию блок не задан и конечно не равен максимальному блоку
}
  
/*
	Функция инициализации блока 
*/
func (b *block) init_block() {
	b.id = len(blocks)   	// номер по порядку создания
	b.size = 0                	// стартовый размер блока равен 0
	b.start = nil     		// стартовая точка неизвестна
	
	blocks = append(blocks, b)
}

/*
	Функция добавления ячейки в состав блока

	Если ячейка уже принадлежит другому блоку, то надо проверить дополнительное условие:
	переход ячейки в новый блок возможен только вместе со всеми остальными ячейками старого блока
	и это делается только если номер нового блока меньше чем номер старого (то есть он создан раньше)
*/
func (b *block) add_cell(c *cell) (result bool, res_desc string) {

//	fmt.Println("Добавляем ячейку в блок", b.id) 

	check_max := func () {

		// если при соединении блоков образовался самый большой блок на поле, 
		// то это надо зафиксировать в глобальных переменных программы

//		fmt.Println("Сравниваем размер максимального блока и текущего блока", max_block_size, b.size) 

		if b.size > max_block_size {
//			fmt.Println("Размер максимального блока увеличивается до ", b.size) 
        	        max_block_id = b.id
			max_block_size = b.size
		}  
	}

   	if c.block_id >= 0  {
		if c.block_id < b.id {
	   		result = false
			res_desc = "Ячейка уже входит в более старый блок"
			return
		} else { // если же ячейка относится к более молодому блоку, то надо весь блок присоединить к текущему
 			b2 := blocks[c.block_id]
//			fmt.Println("Блок поглощается более старым блоком", b2.id) 
			b2.parent_id = b.id
			
			for _, c2 := range b2.cells {
				c2.block_id = b.id
			}

			b.cells = append(b.cells, b2.cells...)
			b.size = len(b.cells)

//			fmt.Println("Размер блока стал", b.size) 
		
			
			check_max()
			res_desc = "Ячейка добавлена в блок вместе с остальными ячейками блока"
			
			return
		}
  	}
	
	result = true
	res_desc = ""

   	if (c.block_id == b.id) {
		res_desc = "Ячейка уже входит в этот блок"
		return
  	}

 	// добавление независимой ячейки в текущий блок
//	fmt.Println("Число ячеек в блоке =", len(b.cells)) 
  	b.cells = append(b.cells, c)    	// ячейка включается в состав этого блока
//	fmt.Println("Новое число ячеек в блоке =", len(b.cells)) 
	b.size = len(b.cells)	  	// размер текущего блока увеличивается на 1
//	fmt.Println("Размер блока =", b.size) 

	if b.start == nil {
 		b.start = c             	// назначается стартовая точка 
	}
	
	c.block_id = b.id			// ячейки приписывается к блоку

 	check_max()

	res_desc = "Ячейка добавлена в блок"
	return
}

/*
	Функция вывода матрицы на экран
*/
func printColors() {
	for i := 0; i < row_count; i++ {
		for j := 0; j < col_count; j++ {
                        if cells[i][j].block_id == max_block_id {
 				fmt.Print("", " ") 
		        } else {
 				fmt.Print(cells[i][j].color, ",") 
// 				fmt.Print("", " ") 
			}
	
		}
		fmt.Println("") 
	}
}

func printBlockIDs() {
	for i := 0; i < row_count; i++ {
		for j := 0; j < col_count; j++ {
                        if cells[i][j].block_id == max_block_id {
 				fmt.Print("*", ",") 
		        } else {
 				fmt.Print(cells[i][j].block_id, ",") 
			}
	
		}
		fmt.Println("") 
	}
}

/*
	Функция вывода блоков на экран
*/
func printBlocks() {
	for i := 0; i < len(blocks); i++ {
		fmt.Println("Block ", i, *blocks[i]) 
	}
}


/*
	Основная функция приложения
*/
func main() {

	var i int = 0
	var j int = 0


	r := rand.New(rand.NewSource(10))

        // инициализируем ячейки в стартовое состояние
	for i = 0; i < row_count; i++ {
		for j = 0; j < col_count; j++ {
			c := new(cell)
			c.init_cell(i, j, 1+r.Intn(color_count))
			cells[i][j] = c
		}
	}

	// вывод первичной матрицы на экран
	printColors()
	fmt.Println("---------------") 
//	printBlockIDs()
//	fmt.Println("---------------") 
//	printBlocks()



        // основной цикл подсчета блоков
	for i = 0; i < row_count; i++ {
//		var s string
//		fmt.Scanln(&s)
//		fmt.Println("Обработка строки: ", i) 

		for j = 0; j < col_count; j++ {
//			fmt.Println("Обработка ячейки: ", i, j) 

 			// если это стартовая ячейка, то создаём первый блок и добавляем в него эту ячейку
			if i == 0 && j == 0 {
//				fmt.Println("Создаём первый блок!") 
				var b block
				b.init_block()
				b.add_cell(cells[i][j])
//				fmt.Println("После добавления в блоке стало ячеек", len(b.cells)) 
				cells[i][j].block_id = b.id
				continue
			}

			// соединение с предыдущим рядом по вертикали, начиная со второго ряда
			if i > 0 && cells[i][j].color == cells[i-1][j].color && cells[i][j].block_id < 0 {
//				fmt.Println("Присоединяем ячейку к верхней") 
				b := blocks[cells[i-1][j].block_id]
				b.add_cell(cells[i][j])
//				fmt.Println("После добавления в блоке стало ячеек", len(b.cells)) 
				cells[i][j].block_id = b.id
			}

                        // соединение с предыдущей ячейкой по горизонтали в случае если текущая ячейка еще не присоединена к блоку
			if j > 0 && cells[i][j].color == cells[i][j-1].color && cells[i][j].block_id < 0 {
//				fmt.Println("Присоединяем ячейку к левой") 
				b := blocks[cells[i][j-1].block_id]
				b.add_cell(cells[i][j])
//				fmt.Println("После добавления в блоке стало ячеек", len(b.cells)) 
				cells[i][j].block_id = b.id
			}

                        // соединение с предыдущей ячейкой по горизонтали в случае если текущая ячейка УЖЕ присоединена к блоку
			if j > 0 && cells[i][j].color == cells[i][j-1].color && cells[i][j].block_id >= 0 && cells[i][j].block_id != cells[i][j-1].block_id {
//				fmt.Println("Сравниваем ячейку к левой") 
  				if cells[i][j].block_id < cells[i][j-1].block_id {
//					fmt.Println("Текущая ячейка поглощает левую") 
					b := blocks[cells[i][j].block_id]
					b.add_cell(cells[i][j-1])
//					fmt.Println("После добавления в блоке стало ячеек", len(b.cells)) 
					cells[i][j-1].block_id = b.id
				} else {
//					fmt.Println("Левая ячейка поглощает текущую") 
					b := blocks[cells[i][j-1].block_id]
					b.add_cell(cells[i][j])
//					fmt.Println("После добавления в блоке стало ячеек", len(b.cells)) 
					cells[i][j].block_id = b.id
				}
				continue
			}


 			// если ячейка не присоединилась ни к слева ни с верху, то создаём новый блок и добавляем в него эту ячейку
			if cells[i][j].block_id < 0 {
//				fmt.Println("Создаём новый блок!") 
				var b block
				b.init_block()
				b.add_cell(cells[i][j])
				cells[i][j].block_id = b.id
			}

		}
	}


//	fmt.Println("---------------") 
	fmt.Println("Max block ID is", max_block_id)
	fmt.Println("Max block size is", max_block_size)
	fmt.Println("----------------------") 

	// вывод матрицы на экран - с выделением максимального блока звездочками
	printColors()
	fmt.Println("---------------") 
//	printBlockIDs()
//	fmt.Println("---------------") 
//	printBlocks()

        // ожидание Enter для выхода из программы
	fmt.Print("Press Enter to exit")
	var s string
	fmt.Scanln(&s)

}
