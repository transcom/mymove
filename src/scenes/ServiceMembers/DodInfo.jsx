import { pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { updateServiceMember, loadServiceMember } from './ducks';
import { reduxifyForm } from 'shared/JsonSchemaForm';
import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';

// todo: add branch (once can get yaml anchors to work)
const subsetOfFields = ['edipi', 'social_security_number', 'rank'];
const uiSchema = {
  title: 'Create your profile',
  order: subsetOfFields,
  requiredFields: subsetOfFields,
  definitions: {
    MilitaryBranch: {
      type: 'string',
      enum: ['ARMY', 'NAVY', 'MARINES', 'AIRFORCE', 'COASTGUARD'],
    },
  },
};

const formName = 'service_member_dod_info';
const CurrentForm = reduxifyForm(formName);

export class ContactInfo extends Component {
  componentDidMount() {
    this.props.loadServiceMember(this.props.match.params.serviceMemberId);
  }

  handleSubmit = () => {
    // const pendingValues = this.props.formData.values;
    // if (pendingValues) {
    //   const patch = pick(pendingValues, subsetOfFields);
    //   this.props.updateServiceMember(patch);
    // }
  };

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentServiceMember,
    } = this.props;
    const isValid = this.refs.currentForm && this.refs.currentForm.valid;
    const isDirty = this.refs.currentForm && this.refs.currentForm.dirty;
    const initialValues = currentServiceMember
      ? pick(currentServiceMember, subsetOfFields)
      : null;
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={isValid}
        pageIsDirty={isDirty}
        hasSucceeded={hasSubmitSuccess}
        error={error}
      >
        <CurrentForm
          ref="currentForm"
          className={formName}
          handleSubmit={no_op}
          schema={this.props.schema}
          uiSchema={uiSchema}
          showSubmit={false}
          initialValues={initialValues}
        />
      </WizardPage>
    );
  }
}
ContactInfo.propTypes = {
  schema: PropTypes.object.isRequired,
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
  const props = {
    userEmail: state.user.email,
    schema: {},
    formData: state.form[formName],
    ...state.serviceMember,
  };
  if (state.swagger.spec) {
    props.schema = state.swagger.spec.definitions.CreateServiceMemberPayload;
  }
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(ContactInfo);
