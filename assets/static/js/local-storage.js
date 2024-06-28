

/**
 * Get the value from localStorage if the key matches
 * @param key localStorage value
 * @param defaultValue if not found return this as the value
 * @returns matching key, or default value
 */
const getValue = (key, defaultValue = {}) => {
    try {
        // read value from local storage
        const item = localStorage.getItem(key);
        return item ? JSON.parse(item) : defaultValue;
    } catch (error) {
        console.log(error);
        return defaultValue;
    }
}

/**
 * Set a value in localStorage
 * @param key name of the value
 * @param value data to store with the key
 */
const setValue = (key, value) => {
    try {
        window.localStorage.setItem(key, JSON.stringify(value));
    } catch (error) {
        console.log(error);
    }
}

// <div x-data="{textToCopy: '{{` $textToCopy `}}'}">
/**
 * appendLocalStorageHashes saves the hashes object to localstorage, appending
 * to existing data
 * @param resp json object
 */
const appendLocalStorageHashes = (resp) => {
    const LSObject = {
        short_url: resp.short_url,
        hash: resp.link.hash,
        long_url: resp.link.original_url
    }
    const allLS = getValue("hashes", '') || [];
    allLS.push(LSObject)
    setValue("hashes", allLS)
    console.log('set values', getValue("hashes"))
}
