$('form').on('submit', function main(e){
    e.preventDefault();
    let a = $('#a').val()
    let b = $('#b').val()
    let c = $('#c').val()
    let error = false 
    
    if (!a) a = 1
    if (!/^-?[0-9]+$/.test(a)) {
        $('#a').css('border-color', 'red')
        error = true
    }

    if (!b) b = 1
    if (!/^-?[0-9]+$/.test(b)) {
        $('#b').css('border-color', 'red')
        error = true
    }
    if (!c) c = 0
    if (!/^-?[0-9]+$/.test(c)) {
        $('#c').css('border-color', 'red')
        error = true
    }
    if (error){}
    else {
        console.log(a, b, c)
        $.ajax(
            {
            url: "/api/solve",
            type: "POST",
        
            data: { a: a,
                    b: b,
                    c: c
                },
            success: function (result) {
                var returnedData = JSON.parse(result);
                if (returnedData.IsValid === true){
                    if (returnedData.Delta == 0){
                        $('#delta').text(returnedData.Delta)
                        $('#x1').text(returnedData.X1.toFixed(2))
                        $('#x2').html('x<sub>1</sub>')
                    }
                    else {
                        $('#delta').text(returnedData.Delta)
                        $('#x1').text(returnedData.X1.toFixed(2))
                        $('#x2').text(returnedData.X2.toFixed(2))
                    }
                }
                else {
                    $('#delta').text(returnedData.Delta)
                    $('#x1').text('Functia are valori imaginare')
                    $('#x2').text('Functia are valori imaginare')
                }
            }
        });  
    }
       
})