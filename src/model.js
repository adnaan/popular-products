import { getProducts, updateVote } from './services';

export default {
    namespace: 'products',
    state: {
        list: [],
        loading: false,
    },
    subscriptions: {
        init({ dispatch }: { dispatch: Function }) {
            dispatch({ type: 'query' });
        },
    },
    effects: { * query(_, { call, put }: { call: Function, put: Function }) {
            const { success, data } = yield call(getProducts);
            if (success) {
                yield put({
                    type: 'querySuccess',
                    products: data,
                });
            }
        },
        * vote({ id }: { id: number }, { call, put }: { call: Function, put: Function }) {
            const { success } = yield call(updateVote, id);
            if (success) {
                yield put({
                    type: 'voteSuccess',
                    id,
                });
            }
        },
    },
    reducers: {
        query(state) {
            return {...state, loading: true, };
        },
        querySuccess(state, { products }: { products: Array < any > }) {
            return {...state, loading: false, list: products };
        },
        vote(state) {
            return {...state, loading: true };
        },
        voteSuccess(state, { id }: { id: number }) {
            const newList = state.list.map(product => {
                if (product.ID === id) {
                    return {...product, vote: product.vote + 1 };
                } else {
                    return product;
                }
            });
            return {...state, list: newList, loading: false };
        },
    },
}
