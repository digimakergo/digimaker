$(document).ready(function(){
    $('.icon-toggle').click(function(){
        var children = $(this).parent().nextAll( 'ul' );
        if( $(this).hasClass( 'open' ) ){
            children.hide();
            $(this).removeClass('open');
            $(this).addClass('closed');

        }else{
            $(this).removeClass('closed');
            $(this).addClass('open');
            children.show();
        }

    })

});


var dm = {
    uploadFile: function(ele){
        var data = new FormData();
        if( $(ele)[0].files.length > 0 )
        {
            var file = $(ele)[0];
            data.append( "file", file.files[0] );
            var fieldName = $(ele).data("field");
            var clientFileName = file.files[0].name;
            clientFileName = clientFileName.substring(0, clientFileName.lastIndexOf("."));
            $.ajax({
                url: '/api/util/uploadimage',
                data: data,
                cache: false,
                contentType: false,
                processData: false,
                method: 'POST',
                type: 'POST', // For jQuery < 1.9
            }).done(function(data){
                $("[name='"+fieldName+"']").val( data );
                $("[name='title']").val(clientFileName);
            }).fail(function(jqXHR, textStatus) {
                alert( "error: " + jqXHR.responseText );
            });
        }
    }

}
