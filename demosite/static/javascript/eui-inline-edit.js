$(document).ready(function(){
    $('.block, .full').mouseover(function(){
        var editor = $(this).find('>.inline-edit');
        editor.show();
        if( $(this).hasClass( 'block' ) ){
            $('.full > .inline-edit').hide();
        }
    });

    $('.block, .full').mouseout(function(){
        $(this).find('.inline-edit').hide();
    });
});


// var treelist = function(){
//     dialog = $( ".treelist" ).dialog({
//       autoOpen: false,
//       height: 400,
//       width: 350,
//       modal: true,
//       buttons: {
//         "Create an account": function(){},
//         Cancel: function() {
//           dialog.dialog( "close" );
//         }
//       },
//       close: function() {
//         form[ 0 ].reset();
//         allFields.removeClass( "ui-state-error" );
//       }
//     });
// }
