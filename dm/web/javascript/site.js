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
