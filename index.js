import Alpine from 'alpinejs';
import collapse from '@alpinejs/collapse'
import Clipboard from "@ryangjchandler/alpine-clipboard"

Alpine.plugin(collapse)
Alpine.plugin(Clipboard.configure({
    onCopy: () => {
        alert("Copied to Clipboard")
    }
}))
window.Alpine = Alpine;
Alpine.start();
