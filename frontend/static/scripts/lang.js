const urlParams = new URLSearchParams(window.location.search);
const lang = urlParams.get("lang");

if (lang) {
  document.querySelectorAll("a").forEach(link => {
    const href = new URL(link.href);
    href.searchParams.set("lang", lang);
    link.href = href.toString();
  });
}

const logout_btn = document.getElementById("log_out")

logout_btn.addEventListener('click', e => {
  fetch('', { 
    method: "POST",
    headers: {
      'Session': sessionToken
    }
  }) 
})
