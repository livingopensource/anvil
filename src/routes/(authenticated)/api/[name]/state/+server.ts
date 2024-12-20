import { json } from '@sveltejs/kit';
import {fetch, Agent} from 'undici';
import { env } from '$env/dynamic/private';
import * as fs from 'node:fs';
import type { operationResponse } from '$lib/server/incus.types';
import { hash } from '$lib/server/utils.js';

export const PUT = async ({locals, url, params}) => {
    const session = await locals.auth()
    if (session == null) {
        return json(403, {})
    }
    const project = hash(session?.user?.email ?? "") ?? "none";
    const state = url.searchParams.get("state")
    if (state != null) {
        const res = await fetch(`${env.CLUSTER_URL}/1.0/instances/${params.name}/state?project=${project}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            dispatcher: new Agent({
                connect: {
                    cert: fs.readFileSync(env.CERT),
                    key: fs.readFileSync(env.KEY),
                    rejectUnauthorized: false,
                }
            }),
            body: JSON.stringify({
                action: state,
                timeout: 10,
                force: false,
                stateful: false,
            })
        })
    
        if (!res.ok) {
            return json({"error": res.statusText}, {
                status: 400
            })
        }
    
        const data = await res.json() as operationResponse
        if (data.status_code == 100) {
            return new Response(String(data))
        }
    
        return json({"error": `unable to ${state} instance`}, {
            status: 400
        })
    }

    return json({"error": "unknown operation"}, {
        status: 400
    })
}