import React from 'react';
import PropTypes from 'prop-types';
import { reduxForm, Field } from 'redux-form';

import SaveCancelButtons from 'scenes/Review/SaveCancelButtons';
import { withContext } from 'shared/AppContext';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import FileUpload from 'components/FileUpload/FileUpload';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { documentSizeLimitMsg } from 'shared/constants';
import { createModifiedSchemaForOrdersTypesFlag } from 'shared/featureFlags';
import { DutyStationSearchBox } from 'scenes/ServiceMembers/DutyStationSearchBox';
import 'scenes/Review/Review.css';
import profileImage from 'scenes/Review/images/profile.png';
import { ExistingUploadsShape } from 'types';
import { OrdersShape } from 'types/customerShapes';

const editOrdersFormName = 'edit_orders';

const EditOrdersForm = ({
  createUpload,
  onDelete,
  schema,
  handleSubmit,
  submitting,
  valid,
  initialValues,
  existingUploads,
  onUploadComplete,
  filePondEl,
  context,
}) => {
  const showAllOrdersTypes = context.flags.allOrdersTypes;
  const modifiedSchemaForOrdersTypesFlag = createModifiedSchemaForOrdersTypesFlag(schema);

  return (
    <div className="grid-container usa-prose">
      <div className="grid-row">
        <div className="grid-col-12">
          <form onSubmit={handleSubmit}>
            <img src={profileImage} alt="" />
            <h1
              style={{
                display: 'inline-block',
                marginLeft: 10,
                marginBottom: 16,
                marginTop: 20,
              }}
            >
              Orders
            </h1>
            <SectionWrapper>
              <h2>Edit Orders:</h2>
              <SwaggerField
                fieldName="orders_type"
                swagger={showAllOrdersTypes ? schema : modifiedSchemaForOrdersTypesFlag}
                required
              />
              <SwaggerField fieldName="issue_date" swagger={schema} required />
              <SwaggerField fieldName="report_by_date" swagger={schema} required />
              <SwaggerField fieldName="has_dependents" swagger={schema} component={YesNoBoolean} />
              <br />
              <Field name="new_duty_station" component={DutyStationSearchBox} />
              <p>Uploads:</p>
              {existingUploads?.length > 0 && <UploadsTable uploads={existingUploads} onDelete={onDelete} />}
              {initialValues?.uploaded_orders && (
                <div>
                  <p>{documentSizeLimitMsg}</p>
                  <FileUpload
                    ref={filePondEl}
                    createUpload={createUpload}
                    onChange={onUploadComplete}
                    labelIdle={'Drag & drop or <span class="filepond--label-action">click to upload orders</span>'}
                  />
                </div>
              )}
            </SectionWrapper>
            <SaveCancelButtons valid={valid} submitting={submitting} />
          </form>
        </div>
      </div>
    </div>
  );
};

EditOrdersForm.propTypes = {
  context: PropTypes.shape({
    flags: PropTypes.shape({
      allOrdersTypes: PropTypes.bool,
    }).isRequired,
  }).isRequired,
  createUpload: PropTypes.func.isRequired,
  onUploadComplete: PropTypes.func.isRequired,
  onDelete: PropTypes.func.isRequired,
  valid: PropTypes.bool.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  existingUploads: ExistingUploadsShape,
  submitting: PropTypes.bool.isRequired,
  filePondEl: PropTypes.shape({
    current: PropTypes.shape({}),
  }),
  schema: PropTypes.shape({}).isRequired,
  initialValues: OrdersShape.isRequired,
};

EditOrdersForm.defaultProps = {
  existingUploads: [],
  filePondEl: null,
};

export default withContext(
  reduxForm({
    form: editOrdersFormName,
  })(EditOrdersForm),
);
