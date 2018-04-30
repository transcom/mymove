import { pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { updateServiceMember, loadServiceMember } from './ducks';
import { reduxifyForm } from 'shared/JsonSchemaForm';
import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';

const subsetOfFields = ['name', 'email', 'telephone'];

const uiSchema = {
  title: 'Backup Contacts',
  description:
    "If we can't reach you, who can we contact (such as spouse or parent)?",
  order: subsetOfFields,

  requiredFields: ['name', 'email'],
  todos: (
    <ul>
      <li>load/save is not wired up (since backend for this is not done)</li>
      <li>
        leaving out permissions ui since we are prioritizing getting flow for
        service member done first and current model for permissions is not in
        master yet
      </li>
    </ul>
  ),
};
const formName = 'service_member_backup_contact';
const CurrentForm = reduxifyForm(formName);

export class BackupContact extends Component {
  componentDidMount() {
    //  this.props.loadServiceMember(this.props.match.params.serviceMemberId);
  }

  handleSubmit = () => {
    // const pendingValues = this.props.formData.values;
    // if (pendingValues) {
    //   debugger;
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
    // initialValues has to be null until there are values from the action since only the first values are taken
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
          initialValues={initialValues}
        />
      </WizardPage>
    );
  }
}
BackupContact.propTypes = {
  userEmail: PropTypes.string.isRequired,
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
    schema: {},
    formData: state.form[formName],
    ...state.serviceMember,
  };
  if (state.swagger.spec) {
    props.schema =
      state.swagger.spec.definitions.CreateServiceMemberBackupContactPayload;
  }
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(BackupContact);
