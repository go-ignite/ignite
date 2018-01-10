var Panel = function () {

    var createHandler = function () {
        $('#server-type').on('change', function (e) {
            var methods=[];
            if(this.value == 'SS') {
                // for ss
                methods=ssMethods;
            } else if(this.value=="SSR"){
                // for ssr
                methods=ssrMethods;
            }
            $("#method").empty();
            if(methods.length==0){
                $("#method").append("<option value='-1'>请选择加密方式</option>");
            }else{
                for (i in methods){
                    $("#method").append("<option value='"+methods[i]+"'>"+methods[i]+"</option>");
                } 
            }
        });

        $('#create-btn').on('click', function (e) {
            e.preventDefault();

            //1. Hide create-btn.
            $('#form-title').css('display', 'none');
            $('#create-btn').css('display', 'none');
            $('#form-method').css('display', 'none');
            $('#form-type').css('display', 'none');

            //2. Show loading icon.
            $('.boxLoading').fadeIn(500);

            //3. Send create SS service request & show account info panel.
            var form = $('#create-form');
            var url = form.attr("action");
            $.post(url, form.serialize(), function (resp) {
                if (resp.success) {
                    $('#host').val(resp.data.host);
                    $('#port').val(resp.data.servicePort);
                    $('#pwd').val(resp.data.servicePwd);
                    $('#encrypt').val($('#method').val());
                    $('#types').val($('#server-type').val());

                    $('#package-limit').html(resp.data.packageLimit+'<up>GB</up>');
                    $('#package-used').html('0<up>GB</up>');
                    $('#package-left').html(resp.data.packageLimit+'<up>GB</up>');
                    $('.progressbar').attr('data-perc', '0');
                    $('#service-status').html('<div class="led green"></div><span>运行中</span>');

                    $('.boxLoading').css('display', 'none');
                    $('.infobox').fadeIn(1500);
                } else {
                    //Create SS service failed
                    toastr.warning(data.message);
                    return false;
                }
            }, "json");
        });
    };

    return {
        //main function to initiate the module
        init: function () {
            createHandler();
        }
    };
}();