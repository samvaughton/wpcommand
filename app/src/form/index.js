class Form {

    constructor(formType, fields, existingObject, fieldsTypeMap) {
        this.formType = formType;
        this.fields = fields;
        this.fieldsTypeMap = fieldsTypeMap || {};
        this.existingObject = existingObject;
        this.submitted = false;
        this.errorMessage = '';

        this.init();
    }

    initWithNewData(object) {
        this.existingObject = object;
        this.init();
    }

    init() {
        this.pristine = {};
        this.current = {};

        this.fields.forEach(field => {
            let value = '';

            if (this.fieldsTypeMap[field] !== undefined) {
                const type = this.fieldsTypeMap[field] !== undefined;

                if (type === 'bool') {
                    value = false;
                }
            }

            if (
                this.formType === 'UPDATE' &&
                this.existingObject !== undefined &&
                this.existingObject !== null &&
                this.existingObject[field] !== undefined
            ) {
                value = this.existingObject[field];
            }

            this.pristine[field] = {
                value: value,
                error: ''
            };
        });

        this.reset(); // clear all errors and assign current
    }

    getValuesFromObj(extra) {
        let values = {};

        for (let key in this.current) {
            values[key] = this.current[key].value;
        }

        console.log(this.current);

        return Object.assign(values, extra);
    }

    reset() {
        this.clearErrors();
        this.current = JSON.parse(JSON.stringify(this.pristine));
    }

    clearErrors() {
        this.submitted = false;
        this.errorMessage = '';
        for (let key in this.current) {
            this.current[key].error = '';
        }
    }

    isValid(item) {
        return this.submitted && item.error === '';
    }

    isInvalid(item) {
        return this.submitted && item.error !== '';
    }

    hydrateErrorsFromValidationErrors(errors) {
        for (let key in errors) {
            if (this.current[key] !== undefined) {
                this.current[key].error = errors[key];
            }
        }

        this.submitted = true;
    }

}

export default Form