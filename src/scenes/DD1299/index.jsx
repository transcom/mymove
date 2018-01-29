import React from 'react';
import DD1299Form from './DD1299Form';
import TestForm from './jsonschema';
export default class DD1299 extends React.Component {
  submit = values => {
    // print the form values to the console
    console.log(values);
  };
  render() {
    //return <DD1299Form onSubmit={this.submit} />;
    return <TestForm />;
  }
}
