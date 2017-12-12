function formatSize(length){
    var i = 0, type = ['б','Кб','Мб','Гб','Тб','Пб'];
    while((length / 1000 | 0) && i < type.length - 1) {
        length /= 1024;
        i++;
    }
    document.getElementById('div_for_insert').innerHTML = length.toFixed(2) + ' ' + type[i];
}