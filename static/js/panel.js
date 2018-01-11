var Panel = function () {

    var createHandler = function () {
        $('#extend').on('click', function(e) {
            if($('#extend').hasClass('fa-angle-double-up')) {
                $('#extend').removeClass('fa-angle-double-up');
                $('#extend').addClass('fa-angle-double-down');
                $('#account-detail').slideToggle('slow');
            } else {
                $('#extend').removeClass('fa-angle-double-down');
                $('#extend').addClass('fa-angle-double-up');
                $('#account-detail').slideToggle('slow');
            }
        });

        $('#server-type').on('change', function (e) {
            var methods = [];
            if (this.value == 'SS') {
                // for ss
                methods = ssMethods;
            } else if (this.value == "SSR") {
                // for ssr
                methods = ssrMethods;
            }
            $("#method").empty();
            if (methods.length == 0) {
                $("#method").append("<option value='-1'>请选择加密方式</option>");
            } else {
                for (i in methods) {
                    $("#method").append("<option value='" + methods[i] + "'>" + methods[i] + "</option>");
                }
            }
        });

        $('#create-btn').on('click', function (e) {
            e.preventDefault();

            // Show loading
            $('#create-form').css('display', 'none');
            $('.boxLoading').fadeIn(500);

            // Send create SS service request & show account info panel.
            var form = $('#create-form');
            var url = form.attr("action");
            $.post(url, form.serialize(), function (resp) {
                if (resp.success) {
                    $('#host').val(resp.data.host);
                    $('#port').val(resp.data.servicePort);
                    $('#pwd').val(resp.data.servicePwd);
                    $('#encrypt').val($('#method').val());
                    $('#types').val($('#server-type').val());

                    $('#package-limit').html(resp.data.packageLimit + '<up>GB</up>');
                    $('#package-used').html('0<up>GB</up>');
                    $('#package-left').html(resp.data.packageLimit + '<up>GB</up>');
                    $('.progressbar').attr('data-perc', '0');
                    $('#service-status').html('<div class="led green"></div><span>运行中</span>');

                    $('.boxLoading').css('display', 'none');
                    $('.infobox').fadeIn(1500);
                } else {
                    $('.boxLoading').css('display', 'block');
                    $('.boxLoading').fadeOut(500);
                    $('#create-form').fadeIn(500);
                    toastr.warning(resp.message);
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
