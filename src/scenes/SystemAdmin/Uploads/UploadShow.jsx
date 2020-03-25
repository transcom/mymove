import React from 'react';
import { ShowController, ShowView, SimpleShowLayout, TextField, DateField } from 'react-admin';

const UploadShow = props => (
  <ShowController {...props}>
    {controllerProps => (
      <ShowView {...props} {...controllerProps}>
        <SimpleShowLayout>
          {controllerProps.record && controllerProps.record.service_member_id
            ? [
                <TextField source="service_member_id" label="Service Member ID" />,
                <TextField source="service_member_first_name" label="Service Member First Name" />,
                <TextField source="service_member_last_name" label="Service Member Last Name" />,
                <TextField source="service_member_phone" label="Service Member Phone" />,
                <TextField source="service_member_email" label="Service Member Email" />,
              ]
            : [
                <TextField source="office_user_id" label="Office User ID" />,
                <TextField source="office_user_first_name" label="Office User First Name" />,
                <TextField source="office_user_last_name" label="Office User Last Name" />,
                <TextField source="office_user_phone" label="Office User Phone" />,
                <TextField source="office_user_email" label="Office User Email" />,
              ]}
          <TextField source="move_locator" label="Move Locator" />
          <TextField source="upload.filename" label="Upload Filename" />
          <TextField source="upload.size" label="Upload Size" />
          <TextField source="upload.content_type" label="Upload Content Type" />
          <DateField source="upload.created_at" showTime label="Created At" />
        </SimpleShowLayout>
      </ShowView>
    )}
  </ShowController>
);

export default UploadShow;
