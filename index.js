import Alpine from 'alpinejs';
import Clipboard from "@ryangjchandler/alpine-clipboard"

Alpine.plugin(Clipboard.configure({
    onCopy: () => {
        alert("Copied to Clipboard")
    }
}))
window.Alpine = Alpine;
Alpine.start();
