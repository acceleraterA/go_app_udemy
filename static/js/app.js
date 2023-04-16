function Prompt() {
    let toast = function(c) {
        const {
            msg = "",
                icon = "success",
                position = "top-end",

        } = c;
        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }
    let success = function(c) {
        const {
            msg = "",
                title = "",
                footer = "",
        } = c;

        Swal.fire({
            icon: "success",
            title: title,
            text: msg,
            footer: footer,
        })
    }
    let error = function(c) {
        const {
            msg = "",
                title = "",
                footer = "",
        } = c;

        Swal.fire({
            icon: "error",
            title: title,
            text: msg,
            footer: footer,
        })
    }

    async function custom(c) {
        const {
            icon = "",
                msg = "",
                title = "",
                showConfirmButton = true,
                showCancelButton = true,
        } = c;

        const {
            value: result
        } = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: showCancelButton,
            showConfirmButton: showConfirmButton,
            //before open 
            willOpen: () => {
                //add range pick canlendar https://github.com/mymth/vanillajs-datepicker
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            //when it's open, remove disable attribute
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            },
            //return [start, end]
            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value
                ]
            }

        })

        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.dismiss !== "") {
                    if (c.callback !== undefined) {
                        c.callback(result);
                    }
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false);
            }
        }

    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,

    }

}