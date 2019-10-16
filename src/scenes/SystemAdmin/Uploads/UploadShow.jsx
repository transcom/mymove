import React from 'react';
import { Component } from 'react';
import { Show, SimpleShowLayout, TextField, DateField } from 'react-admin';

class UploadShow extends Component {
  render(props) {
    return (
      <Show {...this.props}>
        <SimpleShowLayout>
          <TextField source="id" />
          <TextField source="move_locator" />
          <TextField source="upload.filename" />
          <TextField source="upload.size" />
          <TextField source="upload.content_type" />
          <DateField source="upload.created_at" showTime />
        </SimpleShowLayout>
      </Show>
    );
  }
}

export default UploadShow;
