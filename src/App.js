/*@flow*/
import React from 'react';
import ProductList from './components/ProductList'
import { connect } from 'react-redux';

const App = connect(({ products }) => ({
    products
}))(function(props) {
    return (
        <div >
      <h2>Popular Products</h2>
      <ProductList
        data={props.products.list}
        loading={props.products.loading}
        dispatch={props.dispatch}
      />
    </div>
    );
});

export default App;
