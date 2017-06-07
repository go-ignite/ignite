var Panel = function () {

    var createHandler = function () {
        $('#create-btn').on('click', function (e) {
            e.preventDefault();

            //1. Hide create-btn.
            $('#form-title').css('display', 'none');
            $('#create-btn').css('display', 'none');

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

                    $('#package-limit').html(resp.data.packageLimit+'<up>GB</up>');
                    $('#package-used').html('0<up>GB</up>');
                    $('#package-left').html(resp.data.packageLimit+'<up>GB</up>');
                    $('.progressbar').attr('data-perc', '0');

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