import { pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import {
  updateServiceMember,
  loadServiceMember,
  indexBackupContacts,
  createBackupContact,
  updateBackupContact,
} from './ducks';
import { reduxifyForm } from 'shared/JsonSchemaForm';
import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';

const subsetOfFields = ['name', 'email', 'telephone', 'permission'];

const uiSchema = {
  title: 'Backup Contacts',
  description:
    "If we can't reach you, who can we contact (such as spouse or parent)?",
  order: subsetOfFields,

  requiredFields: ['name', 'email', 'permission'],
  todos: (
    <ul>
      <li>Make it a radio button, not a chooser</li>
      <li> load the created.</li>
      <li> setup the second thing</li>
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
    this.props.indexBackupContacts(this.props.match.params.serviceMemberId);
  }

  handleSubmit = () => {
    const pendingValues = this.props.formData.values;
    console.log('SUBMITTING', pendingValues);
    if (pendingValues) {
      if (this.props.currentBackupContacts.length > 0) {
        // update existing
        const oldOne = this.props.currentBackupContacts[0];
        this.props.updateBackupContact(oldOne.id, pendingValues);
      } else {
        this.props.createBackupContact(
          this.props.match.params.serviceMemberId,
          pendingValues,
        );
      }
    }
  };

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentServiceMember,
      userEmail,
    } = this.props;
    const isValid = this.refs.currentForm && this.refs.currentForm.valid;
    const isDirty = this.refs.currentForm && this.refs.currentForm.dirty;
    // initialValues has to be null until there are values from the action since only the first values are taken
    var [backup1, backup2] = this.props.currentBackupContacts;
    const firstInitialValues = backup1 ? pick(backup1, subsetOfFields) : null;
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
          initialValues={firstInitialValues}
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
    {
      updateServiceMember,
      loadServiceMember,
      indexBackupContacts,
      createBackupContact,
      updateBackupContact,
    },
    dispatch,
  );
}
function mapStateToProps(state) {
  const props = {
    userEmail: state.user.email,
    currentBackupContacts: state.serviceMember.currentBackupContacts,
    loggedInUser: state.loggedInUser.loggedInUser,
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
