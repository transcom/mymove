import React, { Component } from 'react';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';
export default class Creator extends Component {
  state = { showForm: false };
  openForm = () => {
    this.setState({ showForm: true });
  };
  closeForm = () => {
    this.setState({ showForm: false });
  };
  render() {
    if (this.state.showForm)
      return (
        <div className="accessorial-panel-modal">
          <div className="title">Add a request</div>
          <div>Form goes here</div>
          <div className="usa-grid-full ">
            <div className="usa-width-one-half">
              <a onClick={this.closeForm}>Cancel</a>
            </div>
            <div className="usa-width-one-half align-right">
              <button className="button button-secondary">
                Save &amp; Add Another
              </button>&nbsp;&nbsp;&nbsp;&nbsp;
              <button className="button button-primary">
                Save &amp; Close
              </button>
            </div>
          </div>
        </div>
      );
    return (
      <a onClick={this.openForm}>
        <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} />
        Add a request
      </a>
    );
  }
}
