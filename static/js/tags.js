$(document).ready(function() {  
    var element = document.getElementsByClassName("form-control cat").length
    var max = 5;
    var x = element;
    $(".add-more").click(function(){   
        if(x < max) {
            x++;
            var html = $(".copy").html();  
            $(".after-add-more").after(html);  
        }
    });  

    $("body").on("click",".remove",function(){   
        $(this).parents(".control-group").remove(); x--;
    });  
  });  
