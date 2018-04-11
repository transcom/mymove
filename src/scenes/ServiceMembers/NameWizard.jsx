import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { updateServiceMember, loadServiceMember } from './ducks';
import WizardPage from 'shared/WizardPage';
import NameForm from './NameForm';

export class NameWizardPage extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Submit Service Member Name';
    this.props.loadServiceMember(this.props.match.params.serviceMemberId);
  }
  handleSubmit = values => {
    console.log('VLUAES', values);
    const nameRequest = {
      serviceMemberId: this.props.match.params.serviceMemberId,
      createServiceMemberPayload: {
        first_name: values.first_name,
        middle_initial: values.middle_initial,
        last_name: values.last_name,
        suffix: values.suffix,
      },
    };
    this.props.updateServiceMember(nameRequest);
  };

  render() {
    const { pages, pageKey, hasSubmitSuccess, error } = this.props;
    const pageIsValid = this.refs.currentForm && this.refs.currentForm.valid;
    const pageIsDirty = this.refs.currentForm && this.refs.currentForm.isDirty;
    return (
      <WizardPage
        // handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={pageIsValid}
        pageIsDirty={pageIsDirty}
        hasSucceeded={hasSubmitSuccess}
        error={error}
      >
        <NameForm onSubmit={this.handleSubmit} />
      </WizardPage>
    );
  }
}

NameWizardPage.propTypes = {
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { updateServiceMember, loadServiceMember },
    dispatch,
  );
}
function mapStateToProps(state) {
  return { ...state.serviceMember };
}
export default connect(mapStateToProps, mapDispatchToProps)(NameWizardPage);
