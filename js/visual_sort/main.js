// https://mp.weixin.qq.com/s?__biz=MzI1MTIzMzI2MA==&mid=2650577234&idx=1&sn=18101b4ce512f594d5c6e3c70e3a2392

window.onload = function () {
    let nums = []
    for (let i = 0; i < 4; i++) {
        // 生成一个 0 - 179的有序数组
        const arr = [...Array(180).keys()] // Array.keys()可以学一下，很有用
        const res = []
        while (arr.length) {
            // 打乱
            const randomIndex = Math.random() * arr.length - 1
            res.push(arr.splice(randomIndex, 1)[0])
        }
        nums = [...nums, ...res]
    }

    const canvas = document.getElementById('canvas')
    const ctx = canvas.getContext('2d')
    ctx.fillStyle = 'white' // 设置画画的颜色
    ctx.translate(500, 500) // 移动中心点到(500, 500)

    // 单个长方形构造函数
    function Rect(x, y, width, height) {
        this.x = x // 坐标x
        this.y = y // 坐标y
        this.width = width // 长方形的宽
        this.height = height // 长方形的高
    }

    // 单个长方形的渲染函数
    Rect.prototype.draw = function () {
        ctx.beginPath() // 开始画一个
        ctx.fillRect(this.x, this.y, this.width, this.height) // 画一个
        ctx.closePath() // 结束画一个
    }

    const CosandSin = []
    for (let i = 0; i < 360; i++) {
        const jiaodu = i / 180 * Math.PI
        CosandSin.push({ cos: Math.cos(jiaodu), sin: Math.sin(jiaodu) })
    }

    function drawAll0(arr) {
        const rects = [] // 用来存储720个长方形
        for (let i = 0; i < arr.length; i++) {
            const num = arr[i]
            const { cos, sin } = CosandSin[Math.floor(i / 2)] // 一个角画两个
            const x = num * cos // x = ρ * cosθ
            const y = num * sin // y = ρ * sinθ
            rects.push(new Rect(x, y, 5, 3)) // 收集所有长方形
        }
        rects.forEach(rect => rect.draw()) // 遍历渲染
    }

    function drawAll(arr) {
        return new Promise((resolve) => {
            setTimeout(() => {
                ctx.clearRect(-500, -500, 1000, 1000) // 清空画布
                const rects = [] // 用来存储720个长方形
                for (let i = 0; i < arr.length; i++) {
                    const num = arr[i]
                    const { cos, sin } = CosandSin[Math.floor(i / 2)] // 一个角画两个
                    const x = num * cos // x = ρ * cosθ
                    const y = num * sin // y = ρ * sinθ
                    rects.push(new Rect(x, y, 5, 3)) // 收集所有长方形
                }
                rects.forEach(rect => rect.draw()) // 遍历渲染
                resolve('draw success')
            }, 10)
        })
    }

    async function bubbleSort(arr) {
        var len = arr.length;
        for (var i = 0; i < len; i++) {
            for (var j = 0; j < len - 1 - i; j++) {
                if (arr[j] > arr[j + 1]) {        //相邻元素两两对比
                    var temp = arr[j + 1];        //元素交换
                    arr[j + 1] = arr[j];
                    arr[j] = temp;
                }
            }
            await drawAll(arr) // 一边排序一边重新画
        }
        return arr;
    }

    async function selectionSort(arr) {
        var len = arr.length;
        var minIndex, temp;
        for (var i = 0; i < len - 1; i++) {
            minIndex = i;
            for (var j = i + 1; j < len; j++) {
                if (arr[j] < arr[minIndex]) {     //寻找最小的数
                    minIndex = j;                 //将最小数的索引保存
                }
            }
            temp = arr[i];
            arr[i] = arr[minIndex];
            arr[minIndex] = temp;
            await drawAll(arr)
        }
        return arr;
    }

    async function insertionSort(arr) {
        if (Object.prototype.toString.call(arr).slice(8, -1) === 'Array') {
            for (var i = 1; i < arr.length; i++) {
                var key = arr[i];
                var j = i - 1;
                while (j >= 0 && arr[j] > key) {
                    arr[j + 1] = arr[j];
                    j--;
                }
                arr[j + 1] = key;
                await drawAll(arr)
            }
            return arr;
        } else {
            return 'arr is not an Array!';
        }
    }

    async function heapSort(array) {
        if (Object.prototype.toString.call(array).slice(8, -1) === 'Array') {
            //建堆
            var heapSize = array.length, temp;
            for (var i = Math.floor(heapSize / 2) - 1; i >= 0; i--) {
                heapify(array, i, heapSize);
                await drawAll(array)
            }

            //堆排序
            for (var j = heapSize - 1; j >= 1; j--) {
                temp = array[0];
                array[0] = array[j];
                array[j] = temp;
                heapify(array, 0, --heapSize);
                await drawAll(array)
            }
            return array;
        } else {
            return 'array is not an Array!';
        }
    }
    function heapify(arr, x, len) {
        if (Object.prototype.toString.call(arr).slice(8, -1) === 'Array' && typeof x === 'number') {
            var l = 2 * x + 1, r = 2 * x + 2, largest = x, temp;
            if (l < len && arr[l] > arr[largest]) {
                largest = l;
            }
            if (r < len && arr[r] > arr[largest]) {
                largest = r;
            }
            if (largest != x) {
                temp = arr[x];
                arr[x] = arr[largest];
                arr[largest] = temp;
                heapify(arr, largest, len);
            }
        } else {
            return 'arr is not an Array or x is not a number!';
        }
    }

    async function quickSort(array, left, right) {
        drawAll(nums)
        if (Object.prototype.toString.call(array).slice(8, -1) === 'Array' && typeof left === 'number' && typeof right === 'number') {
            if (left < right) {
                var x = array[right], i = left - 1, temp;
                for (var j = left; j <= right; j++) {
                    if (array[j] <= x) {
                        i++;
                        temp = array[i];
                        array[i] = array[j];
                        array[j] = temp;
                    }
                }
                await drawAll(nums)
                await quickSort(array, left, i - 1);
                await quickSort(array, i + 1, right);
                await drawAll(nums)
            }
            return array;
        } else {
            return 'array is not an Array or left or right is not a number!';
        }
    }

    async function radixSort(arr, maxDigit) {
        var mod = 10;
        var dev = 1;
        var counter = [];
        for (var i = 0; i < maxDigit; i++, dev *= 10, mod *= 10) {
            for (var j = 0; j < arr.length; j++) {
                var bucket = parseInt((arr[j] % mod) / dev);
                if (counter[bucket] == null) {
                    counter[bucket] = [];
                }
                counter[bucket].push(arr[j]);
            }
            var pos = 0;
            for (var j = 0; j < counter.length; j++) {
                var value = null;
                if (counter[j] != null) {
                    while ((value = counter[j].shift()) != null) {
                        arr[pos++] = value;
                        await drawAll(arr)
                    }
                }
            }
        }
        return arr;
    }

    async function shellSort(arr) {
        var len = arr.length,
            temp,
            gap = 1;
        while (gap < len / 5) {          //动态定义间隔序列
            gap = gap * 5 + 1;
        }
        for (gap; gap > 0; gap = Math.floor(gap / 5)) {
            for (var i = gap; i < len; i++) {
                temp = arr[i];
                for (var j = i - gap; j >= 0 && arr[j] > temp; j -= gap) {
                    arr[j + gap] = arr[j];
                }
                arr[j + gap] = temp;
                await drawAll(arr)
            }
        }
        return arr;
    }

    drawAll(nums) // 执行渲染函数

    document.getElementById('btn').onclick = function () {
        // bubbleSort(nums)
        // selectionSort(nums)
        // insertionSort(nums)
        // heapSort(nums)
        // quickSort(nums, 0, nums.length - 1)
        // radixSort(nums, 3)
        shellSort(nums)
    }
}
