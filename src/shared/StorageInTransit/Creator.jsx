import React, { Component } from 'react';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';
import PropTypes from 'prop-types';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
export class Creator extends Component {
  state = { showForm: false, closeOnSubmit: true };

  render() {
    if (this.state.showForm)
      return {
        /* TODO: Render form here */
      };
    return (
      <div className="add-request">
        <a onClick={this.openForm}>
          <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} />
          Request SIT
        </a>
      </div>
    );
  }
}
Creator.propTypes = {
  SITRequests: PropTypes.array,
};

function mapStateToProps(state) {
  return {
    /*
    formEnabled: isValid(SITRequestFormName)(state) && !isSubmitting(SITRequestFormName)(state),
    hasSubmitSucceeded: hasSubmitSucceeded(SITRequestFormName)(state),
    */
  };
}

function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators(
    {
      /*
      submitForm: () => submit(SITRequestFormName),
      clearForm: () => reset(SITRequestFormName),
      */
    },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(Creator);
