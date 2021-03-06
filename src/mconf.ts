// This code is a part of MagicCap which is a MPL-2.0 licensed project.
// Copyright (C) Jake Gealer <jake@gealer.email> 2019.

// Imports go here.
import * as i18n from "./i18n"
import config from "./config"
import { importedUploaders } from "./uploaders"

/**
 * Parses the *.mconf file.
 * @param {Object} data - The JSON parsed data.
 * @returns - All the configuration items that are being changed.
 */
export async function parse(data: {
    version: number;
    config_items: undefined | object;
}) {
    const version = data.version as number
    if (version !== 1) {
        const wrongVerErr = await i18n.getPoPhrase("This version of MagicCap cannot read the config file given.", "mconf")
        throw new Error(wrongVerErr)
    }
    if (data.config_items === undefined || typeof data.config_items !== "object") {
        const cantParseErr = await i18n.getPoPhrase("MagicCap couldn't parse the config file.", "mconf")
        throw new Error(cantParseErr)
    }
    return data.config_items
}


/**
 * Gets the values of a object.
 *
 * @param {any} item - The item you want values from.
 * @returns The values.
 */
function values(item: any) {
    const x: any[] = []
    for (const i in item) {
        const y = item[i] as any
        x.push(y)
    }
    return x
}

/**
 * Handles making a new *.mconf's file contents.
 * @returns The parsed mconf file.
 */
export function newConfig() {
    const options: any = {}
    for (const uploader of values(importedUploaders)) {
        for (const option of values(uploader.config_options)) {
            if (config.o[option.value] !== undefined) {
                options[option.value] = config.o[option.value]
            }
        }
    }
    return {
        version: 1,
        config_items: options,
    }
}
