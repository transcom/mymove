import React from 'react';

import { Component, Fragment } from 'react';
import { Redirect } from 'react-router-dom';

export class GexSearch extends Component {
  state = { ...this.initialState };

  get initialState() {
    return {
      showUpload: false,
    };
  }

  redirectToShowUpload = () => {
    this.setState({ showUpload: true });
  };

  sendRequest = (values) => {
    console.log('sendRequest...', values);
    // this.props.sendGexRequest(values).then(response => {
    //   this.setState({
    //     response: get(response, 'payload.gex_response', 'No payload'),
    //   });
    // });
  };

  render() {
    if (!this.state.showUpload) {
      return (
        <Fragment>
          <span>Send to Gex</span>
          <form onSubmit={this.sendRequest}>
            <input name="reportName" component="input" type="text" />
            <button type="submit">Send</button>
          </form>
        </Fragment>
      );
    } else {
      return <Redirect to={`/uploads/${this.state.uploadID}/show`} />;
    }
  }
}

export default GexSearch;
