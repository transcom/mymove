import React, { Component } from 'react';
import NameForm from './NameForm';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import { loadServiceMember } from './ducks';

class SMName extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Service Member Name';
  }
  render() {
    return (
      <div>
        <NameForm onSubmit={() => {}} />
      </div>
    );
  }
}

SMName.propTypes = {
  currentServiceMember: PropTypes.object,
  currentForm: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadServiceMember }, dispatch);
}

function mapStateToProps(state) {
  return { ...state.serviceMember, currentForm: state.form };
}

export default connect(mapStateToProps, mapDispatchToProps)(SMName);
