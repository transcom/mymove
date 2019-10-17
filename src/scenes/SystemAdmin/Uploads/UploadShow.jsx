import React from 'react';
import { Component } from 'react';
import { Show, SimpleShowLayout, TextField, DateField } from 'react-admin';

class UploadShow extends Component {
  render(props) {
    return (
      <Show {...this.props}>
        <SimpleShowLayout>
          <TextField source="id" label="Upload ID" />
          <TextField source="service_member_id" label="Service Member ID" />
          <TextField source="office_user_id" label="Office User ID" />
          <TextField source="office_user_email" label="Office User Email" />
          <TextField source="move_locator" label="Move Locator" />
          <TextField source="upload.filename" label="Upload Filename" />
          <TextField source="upload.size" label="Upload Size" />
          <TextField source="upload.content_type" label="Upload Content Type" />
          <DateField source="upload.created_at" showTime label="Created At" />
        </SimpleShowLayout>
      </Show>
    );
  }
}

export default UploadShow;
