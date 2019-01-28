import React, { Component } from 'react';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';

import './StorageInTransitPanel.css';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

export class Creator extends Component {
  state = {};

  openForm = e => {
    e.preventDefault();
  };

  render() {
    return (
      <div className="add-request storage-in-transit-hr-top">
        <a onClick={this.openForm}>
          <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} />
          Request SIT
        </a>
      </div>
    );
  }
}

Creator.propTypes = {};

function mapStateToProps(state) {
  return {};
}

function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators({}, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(Creator);
