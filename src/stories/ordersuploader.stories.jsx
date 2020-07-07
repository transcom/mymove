import React from 'react';
import { storiesOf } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import * as Yup from 'yup';

import { OrdersUploader } from 'components/OrdersUploader';

// Orders Uploader
storiesOf('Components|Uploaders', module).add('orders uploader', () => (
  <div>
    <OrdersUploader
      createUpload={action('Createupload')}
      deleteUpload={action('deleteUpload')}
      document={Yup.object({ id: '123' })}
      onChange={action('addToState')}
      options={{ labelIdle: 'Drag & drop or <span class="filepond--label-action">click to upload orders</span>' }}
    />
  </div>
));
