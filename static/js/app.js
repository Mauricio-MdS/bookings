function prompt() {
  const custom = async (c) => {
    const { icon = "", msg = "", showConfirmButton = true, title = "" } = c

    const { value: result } = await Swal.fire({
      backdrop: false,
      focusConfirm: false,
      html: msg,
      icon: icon,
      showCancelButton: true,
      showConfirmButton: showConfirmButton,
      title: title,
      willOpen: () => {
        if (c.willOpen !== undefined) {
          c.willOpen()
        }
      },
      preConfirm: () => {
        return [
          document.getElementById("start").value,
          document.getElementById("end").value,
        ]
      },
      didOpen: () => {
        if (c.didOpen !== undefined) {
          c.didOpen()
        }
      },
    })
    if (result) {
      if (result.dismiss !== Swal.DismissReason.cancel) {
        if (result.value !== "") {
          if (c.callback !== undefined) {
            c.callback(result)
          }
        } else {
          c.callback(false)
        }
      } else {
        c.callback(false)
      }
    }
  }

  const error = (c) => {
    const { msg = "", title = "", footer = "" } = c
    Swal.fire({
      icon: "error",
      title: title,
      text: msg,
      footer: footer,
    })
  }

  const success = (c) => {
    const { msg = "", title = "", footer = "" } = c
    Swal.fire({
      icon: "success",
      title: title,
      text: msg,
      footer: footer,
    })
  }

  const toast = (c) => {
    const { msg = "", icon = "success", position = "top-end" } = c

    const Toast = Swal.mixin({
      toast: true,
      title: msg,
      position: position,
      icon: icon,
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
      didOpen: (toast) => {
        toast.onmouseenter = Swal.stopTimer
        toast.onmouseleave = Swal.resumeTimer
      },
    })
    Toast.fire({})
  }

  return { custom, error, success, toast }
}
