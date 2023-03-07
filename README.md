[TOC]



# JSON Schema

## JSON Schema是什么

看https://json-schema.org/对JSON Schema的定义：`JSON Schema is a declarative language that allows you to annotate and validate JSON documents.`

简单来说，JSON Schema是用来描述一个json，并可以用来校验json的一种声明性语言。

比如有如下一个json：

```json
{
  "name": "eachain"
}
```

用JSON Schema来简略地描述是：

```json
{
  "type": "object", // 首先表明以上json必须是个object
  "properties": { // 对object下面的字段进行描述
    "name": { // 对name字段进行描述
      "type": "string" // name字段必须是string类型
    }
  }
}
```

并且可以用该JSON Schema来校验json格式是否合法。（在线校验：https://jsonschemalint.com/）



## JSON Schema解决什么问题

看https://json-schema.org/对JSON Schema的用途说明：`JSON Schema enables the confident and reliable use of the JSON data format.`即通过了JSON Schema校验的json，程序员在后续业务中可以放心地按schema描述使用json。

比如还是以上例子：

```json
{
  "name": "eachain"
}
```

如果通过了以下JSON Schema的校验：

```json
{
  "type": "object",
  "properties": {
    "name": {
      "type": "string"
    }
  },
  "required": [ // 必须存在的字段
    "name"
  ]
}
```

则在业务中，可以放心地写以下代码（以PHP为例）：

```php
<?php
$request = '{"name":"eachain"}';
// JSON Schema validate $request here.
$json = json_decode($request, true);
echo $json['name'] . "\n";
```

因为在JSON Schema中严格要求必须有`name`字段，如果没有`name`字段，在用JSON Schema校验时不能通过，自然也不会走到这段业务逻辑代码。



## JSON Schema常见应用场景

那么JSON Schema一般用来解决哪些问题呢？

- Web服务器参数校验中间件
    - 经过JSON Schema校验过的json，在后续业务逻辑便可以放心地使用JSON Schema中声明过的结构。
- swagger
    - swagger中对于接口json的描述也是基于JSON Schema的，只要写好JSON Schema，swagger相应工具可以自动生成文档，甚至服务端和客户端代码。从而解决了客户端、服务端、文档的一致性问题。
- 前端表单自动生成
    - JSON Schema相当于是json的元数据，有关于json的描述信息，前端可以根据描述自动生成表单，从而保证在提交表单前，数据就是符合Schema规则的。



## JSON Schema优势

- JSON Schema定义了大量关键字，用于校验不同类型的数据，以及数据格式，可以适应绝大多数场景；
- JSON Schema减轻了开发过程中大量的重复的数据校验工作，使开发人员专注于业务逻辑。
- JSON Schema可以自动生成接口文档，甚至代码，减少了一部分开发工作量。



## JSON Schema劣势

- 版本问题，JSON Schema各版本之间有不兼容的情况，一个项目开发过程中，最好只选定一个版本投入使用
    - 比如在`draft-04`中，关键字`exclusiveMinimum`必须为`boolean`值，而在后续版本中改为了`number`值，用来表示某字段的值必须大于多少
- 关键字拼错不提醒，在开发人员写完JSON Schema后，最好用**定制版**JSON Schema Meta先校验一下JSON Schema，这会有一种无限套娃的感觉
    - 比如业务要求必须要有`name`字段，但不小心将`required`关键字拼错了，此时JSON Schema不会给出任何提示，从而造成服务端拿到的数据不符合预期
- 写法太过灵活，一个“茴”字，真的需要这么多种写法吗？
    - 由其在用到组合关键字时，开发以为这样写是对的，但结果并不一定是对的，需要再拿各种数据来验证一遍。此时开发可能会怀疑：为什么要用JSON Schema呢？



**关于定制版JSON Schema Meta**：业务中用到JSON Schema时，最好限制一下写法，不要太灵活了，因此需要一个JSON Schema来校验开发写的JSON Schema，它被称为JSON Schema Meta。这在下一节会提及。

## 应用JSON Schema时的取舍

有时候特性太多，可选择的写法太灵活，并不是一件好事。我们在用JSON Schema时，最好做出一些取舍，限制用法，这对一个工程项目来说大有裨益。

以下，以比较基础的`draft-04`关键字进行说明。



### type

在实际项目中，应该要求必须出现`type`关键字，并且类型只用`string`。（`type`允许用`array`）

当`type`关键字类型为`array`时，说明其它关键字在校验时只有一部分起作用，一部分不起作用。

- 错误示例

    ```json
    {
      "type": ["string", "number"],
      "minLength": 1, // 只有当type为string时生效
      "minimum": 1 // 只有当type为number时生效
    }
    ```

- 正确示例

    ```json
    {
      "oneOf": [
        {
          "type": "string",
          "minLength": 1
        },
        {
          "type": "number",
          "minimum": 1
        }
      ]
    }
    ```



**关于`integer`类型**：我个人认为这是JSON Schema设计上的一个失误，它不应该通过`{"type":"integer"}`的方式指定整数，而应该以一个关键字的方式指定，比如：

```json
{
  "type": "number",
  "isInteger": true // 注意，实际JSON Schema中无该关键字
}
```





### items

JSON很灵活，一个数组中所有元素类型可以完全不同。但在实际项目中，最好保持元素类型相同。因此，`items`关键字最好只用`object`类型。

- 错误示例

    ```json
    {
      "items": [
        {
          "type": "string"
        },
        {
          "type": "number"
        }
      ]
    }
    ```

- 正确示例

    ```json
    {
      "type": "array",
      "items": {
        "type": "string"
      }
    }
    ```

    

在此情况下，`additionalItems`关键字应该被禁用。



### patternProperties

实际项目中应该禁用该关键字，项目中所用的字段应该是确定的，对于陌生字段，要么忽略不解析，要么不允许有多余字段。如果对多个规则相同的字段用正则表达式来写成一个，那更是大错特错。

- 错误示例

    ```json
    {
      "type": "object",
      "patternProperties": {
        "id$": {
          "type": "number"
        }
      }
    }
    ```

- 正确示例

    ```json
    {
      "type": "object",
      "properties": {
        "aid": {
          "type": "number"
        },
        "bid": {
          "type": "number"
        }
      }
    }
    ```

    



