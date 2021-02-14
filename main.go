package main

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	/*imgName := "./00000"
	imgPath := imgName+".png"
	img, err := OpenFileIntoImage(imgPath)
	if err != nil {
		panic(err)
	}

	OpacityAdjustX(img, 1)

	SaveImage("./00200-2.png", img)*/

	//dxSet, dySet := getDirFile("D:\\BaiduNetdiskDownload\\传奇素材\\57件衣服武器\\25\\衣服\\外观")
	dxSet, dySet := getDirFile("./")
	sortDx := dxSet.ToSlice()
	sort.Slice(sortDx, func(i, j int) bool {
		return i > j
	})
	fmt.Print(sortDx)
	sort.Slice(dySet.ToSlice(), func(i, j int) bool {
		return i > j
	})
	minWidth = float64(sortDx[0].(int))

	for _, filePath := range slice {
		file, _ := os.Open(filePath) //打开图片1
		var (
			imgx  image.Image
			errde error
		)
		if imgx, _, errde = image.Decode(file); errde != nil {
			log.Fatal(errde)
			return
		}
		bx := imgx.Bounds()
		newW := fixSizeX(bx.Max.X)

		// 调用resize库进行图片缩放(高度填0，resize.Resize函数中会自动计算缩放图片的宽高比)
		m1 := resize.Resize(uint(newW), 0, imgx, resize.Lanczos3)
		// 将两个图片合成一张
		newWidth := m1.Bounds().Max.X                                                                         //新宽度 = 随意一张图片的宽度
		newHeight := m1.Bounds().Max.Y + m2.Bounds().Max.Y                                                    // 新图片的高度为两张图片高度的和
		newImg := image.NewNRGBA(image.Rect(0, 0, newWidth, newHeight))                                       //创建一个新RGBA图像
		draw.Draw(newImg, newImg.Bounds(), m1, m1.Bounds().Min, draw.Src)                                     //画上第一张缩放后的图片
		draw.Draw(newImg, newImg.Bounds(), m2, m2.Bounds().Min.Sub(image.Pt(0, m1.Bounds().Max.Y)), draw.Src) //画上第二张缩放后的图片（这里需要注意Y值的起始位置）
		// 保存文件
		imgfile, _ := os.Create("003.png")
		defer imgfile.Close()
		png.Encode(imgfile, newImg)
	}

	/*file1, _ := os.Open("00000.PNG") //打开图片1
	file2, _ := os.Open("00001.PNG") //打开图片2
	defer file1.Close()
	defer file2.Close()
	// image.Decode 图片
	var (
		img1, img2 image.Image
		err        error
	)
	if img1, _, err = image.Decode(file1); err != nil {
		log.Fatal(err)
		return
	}
	if img2, _, err = image.Decode(file2); err != nil {
		log.Fatal(err)
		return
	}
	b1 := img1.Bounds()
	b2 := img2.Bounds()
	new1W, new2W := fixSize(b1.Max.X, b2.Max.X)

	// 调用resize库进行图片缩放(高度填0，resize.Resize函数中会自动计算缩放图片的宽高比)
	m1 := resize.Resize(uint(new1W), 0, img1, resize.Lanczos3)
	m2 := resize.Resize(uint(new2W), 0, img2, resize.Lanczos3)

	// 将两个图片合成一张
	newWidth := m1.Bounds().Max.X       //新宽度 = 随意一张图片的宽度
	newHeight := m1.Bounds().Max.Y + m2.Bounds().Max.Y // 新图片的高度为两张图片高度的和
	newImg := image.NewNRGBA(image.Rect(0, 0, newWidth, newHeight)) //创建一个新RGBA图像
	draw.Draw(newImg, newImg.Bounds(), m1, m1.Bounds().Min, draw.Src) //画上第一张缩放后的图片
	draw.Draw(newImg, newImg.Bounds(), m2, m2.Bounds().Min.Sub(image.Pt(0, m1.Bounds().Max.Y)), draw.Src) //画上第二张缩放后的图片（这里需要注意Y值的起始位置）
	// 保存文件
	imgfile, _ := os.Create("003.png")
	defer imgfile.Close()
	png.Encode(imgfile, newImg)*/
}

const MaxWidth float64 = 600

var minWidth float64

func fixSize(img1W, img2W int) (new1W, new2W int) {
	var ( //为了方便计算，将两个图片的宽转为 float64
		img1Width, img2Width = float64(img1W), float64(img2W)
		ratio1, ratio2       float64
	)

	if minWidth > 600 { // 如果最小宽度大于600，那么两张图片都需要进行缩放
		ratio1 = MaxWidth / img1Width // 图片1的缩放比例
		ratio2 = MaxWidth / img2Width // 图片2的缩放比例

		// 原宽度 * 比例 = 新宽度
		return int(img1Width * ratio1), int(img2Width * ratio2)
	}

	// 如果最小宽度小于600，那么需要将较大的图片缩放，使得两张图片的宽度一致
	if minWidth == img1Width {
		ratio2 = minWidth / img2Width // 图片2的缩放比例
		return img1W, int(img2Width * ratio2)
	}

	ratio1 = minWidth / img1Width // 图片1的缩放比例
	return int(img1Width * ratio1), img2W
}

func fixSizeX(img1W int) (new1W int) {
	var ( //为了方便计算，将两个图片的宽转为 float64
		img1Width = float64(img1W)
		ratio1    float64
	)

	ratio1 = minWidth / img1Width // 图片1的缩放比例
	return int(img1Width * ratio1)
}

//至于为啥用RGBA64，因为任性

//输入图像文件路径，返回 *image.RGBA64 结果
func OpenFileIntoImage(fileName string) (*image.RGBA64, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	//自作聪明检查文件类型
	split := strings.Split(fileName, ".")
	if len(split) <= 1 {
		return nil, fmt.Errorf("i don't think selected file is not a normal png or jpeg image: %s", fileName)
	}
	imageType := split[len(split)-1]
	var imgRes *image.RGBA64
	switch imageType {
	case "jpg", "jpeg":
		img, err := jpeg.Decode(f)
		if err != nil {
			return nil, err
		}
		imgRes = ImageTypeToRGBA64(&img)
	case "png":
		img, err := png.Decode(f)
		if err != nil {
			return nil, err
		}
		imgRes = ImageTypeToRGBA64(&img)
	default:
		return nil, fmt.Errorf("this image format is unknown or not supported yet: %v", imageType)
	}
	return imgRes, nil
}

//Image转换为image.RGBA64
func ImageTypeToRGBA64(m *image.Image) *image.RGBA64 {
	bounds := (*m).Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA64(bounds)
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := (*m).At(i, j)
			r, g, b, a := colorRgb.RGBA()
			nR := uint16(r)
			nG := uint16(g)
			nB := uint16(b)
			alpha := uint16(a)
			newRgba.SetRGBA64(i, j, color.RGBA64{R: nR, G: nG, B: nB, A: alpha})
		}
	}
	return newRgba
}

//将输入图像m的透明度变为原来的倍数。若原来为完成全不透明，则percentage = 0.5将变为半透明
func OpacityAdjust(m *image.RGBA64, percentage float64) *image.RGBA64 {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA64(bounds)
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			r, g, b, a := colorRgb.RGBA()
			opacity := uint16(float64(a) * percentage)
			//颜色模型转换，至关重要！
			v := newRgba.ColorModel().Convert(color.NRGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: opacity})
			//Alpha = 0: Full transparent
			rr, gg, bb, aa := v.RGBA()
			newRgba.SetRGBA64(i, j, color.RGBA64{R: uint16(rr), G: uint16(gg), B: uint16(bb), A: uint16(aa)})
		}
	}
	return newRgba
}

var y1 int
var y2 int

//将输入图像m的透明度变为原来的倍数。若原来为完成全不透明，则percentage = 0.5将变为半透明
func OpacityAdjustX(m *image.RGBA64, percentage float64) {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	y1 = dy
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			_, _, _, a := colorRgb.RGBA()
			if a == 0 {
				continue
			} else {
				if y1 >= j {
					fmt.Printf("x1: %d y1: %d  dy: %d a: %d \n", i, y1, j, a)
					y1 = j
				}
				if y2 <= j {
					fmt.Printf("x2: %d y2: %d  dy: %d a: %d \n", i, y2, j, a)
					y2 = j
				}
				/*opacity := uint16(float64(a)*percentage)
				//颜色模型转换，至关重要！
				v := newRgba.ColorModel().Convert(color.NRGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: opacity})
				//Alpha = 0: Full transparent
				rr, gg, bb, aa := v.RGBA()
				newRgba.SetRGBA64(i, j, color.RGBA64{R: uint16(rr), G: uint16(gg), B: uint16(bb), A: uint16(aa)})*/
			}
		}
	}
}

//保存image.RGBA64为png文件
func SaveImage(fileName string, m *image.RGBA64) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()
	fmt.Println(y2)
	x := m.SubImage(image.Rect(0, y1, m.Bounds().Dx(), y2))
	err = png.Encode(f, x)
	if err != nil {
		panic(err)
	}
}

var slice []string

//递归遍历文件目录
func getDirFile(pathName string) (dxSet, dySet mapset.Set) {
	rd, err := ioutil.ReadDir(pathName)
	if err != nil {
		log.Fatalf("Read dir '%s' failed: %v", pathName, err)
	}
	var name, fullName string
	dxSet = mapset.NewSet()
	dySet = mapset.NewSet()
	for _, fileDir := range rd {
		name = fileDir.Name()
		fullName = pathName + "/" + name
		slice = append(slice, fullName)
		dx, dy := OpenFileGetDxDy(fullName)
		dxSet.Add(dx)
		dySet.Add(dy)
	}
	return dxSet, dySet
}

//输入图像文件路径，返回 *image.RGBA64 结果
func OpenFileGetDxDy(fileName string) (int, int) {
	f, err := os.Open(fileName)
	if err != nil {
	}
	defer func() {
		_ = f.Close()
	}()
	//自作聪明检查文件类型
	split := strings.Split(fileName, ".")
	if len(split) <= 1 {
		fmt.Errorf("i don't think selected file is not a normal png or jpeg image: %s", fileName)
	}
	imageType := split[len(split)-1]
	var imgRes *image.RGBA64
	var dx, dy int
	switch imageType {
	case "jpg", "jpeg":
		break
	case "png", "PNG":
		img, err := png.Decode(f)
		if err != nil {

		}
		imgRes = ImageTypeToRGBA64(&img)
		bounds := (*imgRes).Bounds()
		dx = bounds.Dx()
		dy = bounds.Dy()
	default:
		fmt.Errorf("this image format is unknown or not supported yet: %v", imageType)
	}
	return dx, dy
}
