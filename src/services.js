/*@flow*/
export async function getProducts() {
    return fetch(`/api/products`).then(res => res.json())
}

export async function updateVote(id: { id: number }) {
    return fetch('/api/products/vote/' + String(id)).then(res => res.json())
}
