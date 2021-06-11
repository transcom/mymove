import { React, createRef } from 'react';
import { GridContainer, Grid } from '@trussworks/react-uswds';

import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import ScrollToTop from 'components/ScrollToTop';
import FileUpload from 'components/FileUpload/FileUpload';
import { UploadsShape } from 'types/customerShapes';

const AmendOrders = ({ uploads }) => {
  const handleDelete = () => {};
  const handleUpload = () => {};
  const onChange = () => {};
  const handleSave = () => {};
  const handleCancel = () => {};

  const filePondEl = createRef();

  return (
    <GridContainer>
      <ScrollToTop />
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <h1>Orders</h1>
          <p>
            Upload any amended orders here. The office will update your move info to match the new orders. Talk directly
            with your movers to coordinate any changes.
          </p>
        </Grid>
      </Grid>
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <SectionWrapper>
            <h1>Upload your orders</h1>
            <p>PDF, JPG, or PNG only. Maximum file size 25MB. Each page must be clear and legible</p>
            {uploads && uploads.length > 0 && (
              <>
                <br />
                <UploadsTable uploads={uploads} onDelete={handleDelete} />
              </>
            )}
            <div className="uploader-box">
              <FileUpload
                ref={filePondEl}
                createUpload={handleUpload}
                onChange={onChange}
                labelIdle={'Drag & drop or <span class="filepond--label-action">click to upload orders</span>'}
              />
            </div>
            <WizardNavigation editMode disableNext={false} onNextClick={handleSave} onCancelClick={handleCancel} />
          </SectionWrapper>
        </Grid>
      </Grid>
    </GridContainer>
  );
};

AmendOrders.propTypes = {
  uploads: UploadsShape,
};

// Temporary fake data
AmendOrders.defaultProps = {
  uploads: [
    {
      bytes: 9043,
      content_type: 'image/png',
      created_at: '2021-06-08T23:10:59.442Z',
      filename: 'orders.png',
      id: '2e5b2961-ca0c-4820-aa9d-a693c631b068',
      status: 'PROCESSING',
      updated_at: '2021-06-08T23:10:59.442Z',
      url:
        '/storage/user/11a6ca88-cf5a-4cfc-b4db-9a4cd0285bf3/uploads/2e5b2961-ca0c-4820-aa9d-a693c631b068?contentType=image%2Fpng',
    },
    {
      bytes: 9043,
      content_type: 'image/png',
      created_at: '2021-06-08T23:11:31.885Z',
      filename: 'orders.png',
      id: '74455c13-f2d6-47bb-9cd1-e1fb712e3d0f',
      status: 'PROCESSING',
      updated_at: '2021-06-08T23:11:31.885Z',
      url:
        '/storage/user/11a6ca88-cf5a-4cfc-b4db-9a4cd0285bf3/uploads/74455c13-f2d6-47bb-9cd1-e1fb712e3d0f?contentType=image%2Fpng',
    },
  ],
};
export default AmendOrders;
